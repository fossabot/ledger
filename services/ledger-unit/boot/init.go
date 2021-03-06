// Copyright (c) 2016-2020, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boot

import (
	"context"
	"os"

	"github.com/rs/xid"

	"github.com/jancajthaml-openbank/ledger-unit/actor"
	"github.com/jancajthaml-openbank/ledger-unit/config"
	"github.com/jancajthaml-openbank/ledger-unit/logging"
	"github.com/jancajthaml-openbank/ledger-unit/metrics"
	"github.com/jancajthaml-openbank/ledger-unit/model"
	"github.com/jancajthaml-openbank/ledger-unit/utils"

	system "github.com/jancajthaml-openbank/actor-system"
	localfs "github.com/jancajthaml-openbank/local-fs"
)

// Program encapsulate initialized application
type Program struct {
	interrupt chan os.Signal
	cfg       config.Configuration
	daemons   []utils.Daemon
	cancel    context.CancelFunc
}

// Initialize application
func Initialize() Program {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := config.GetConfig()

	logging.SetupLogger(cfg.LogLevel)

	storage := localfs.NewPlaintextStorage(
		cfg.RootStorage,
	)
	metricsDaemon := metrics.NewMetrics(
		ctx,
		cfg.MetricsOutput,
		cfg.Tenant,
		cfg.MetricsRefreshRate,
	)
	actorSystemDaemon := actor.NewActorSystem(
		ctx, cfg.Tenant,
		cfg.LakeHostname,
		&metricsDaemon,
		&storage,
	)
	transactionFinalizerDaemon := actor.NewTransactionFinalizer(
		ctx,
		cfg.TransactionIntegrityScanInterval,
		&metricsDaemon,
		&storage,
		func(transaction model.Transaction) {
			name := "transaction/" + xid.New().String()
			ref, err := actor.NewTransactionActor(&actorSystemDaemon, name)
			if err != nil {
				return
			}
			ref.Tell(
				transaction,
				system.Coordinates{
					Region: actorSystemDaemon.Name,
					Name:   name,
				},
				system.Coordinates{
					Region: actorSystemDaemon.Name,
					Name:   "transaction_finalizer_cron",
				},
			)
		},
	)

	var daemons = make([]utils.Daemon, 0)
	daemons = append(daemons, metricsDaemon)
	daemons = append(daemons, actorSystemDaemon)
	daemons = append(daemons, transactionFinalizerDaemon)

	return Program{
		interrupt: make(chan os.Signal, 1),
		cfg:       cfg,
		daemons:   daemons,
		cancel:    cancel,
	}
}
