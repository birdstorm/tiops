// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"github.com/joomcode/errorx"
	"github.com/pingcap/errors"
	"github.com/pingcap/tiup/pkg/cluster/meta"
	operator "github.com/pingcap/tiup/pkg/cluster/operation"
	"github.com/pingcap/tiup/pkg/cluster/task"
	"github.com/pingcap/tiup/pkg/logger"
	"github.com/pingcap/tiup/pkg/logger/log"
	"github.com/pingcap/tiup/pkg/utils"
	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop <cluster-name>",
		Short: "Stop a TiDB cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			if err := validRoles(gOpt.Roles); err != nil {
				return err
			}

			clusterName := args[0]
			if utils.IsNotExist(meta.ClusterPath(clusterName, meta.MetaFileName)) {
				return errors.Errorf("cannot stop non-exists cluster %s", clusterName)
			}

			logger.EnableAuditLog()
			metadata, err := meta.ClusterMetadata(clusterName)
			if err != nil {
				return err
			}

			t := task.NewBuilder().
				SSHKeySet(
					meta.ClusterPath(clusterName, "ssh", "id_rsa"),
					meta.ClusterPath(clusterName, "ssh", "id_rsa.pub")).
				ClusterSSH(metadata.Topology, metadata.User, gOpt.SSHTimeout).
				ClusterOperate(metadata.Topology, operator.StopOperation, gOpt).
				Build()

			if err := t.Execute(task.NewContext()); err != nil {
				if errorx.Cast(err) != nil {
					// FIXME: Map possible task errors and give suggestions.
					return err
				}
				return errors.Trace(err)
			}

			log.Infof("Stopped cluster `%s` successfully", clusterName)

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&gOpt.Roles, "role", "R", nil, "Only stop specified roles")
	cmd.Flags().StringSliceVarP(&gOpt.Nodes, "node", "N", nil, "Only stop specified nodes")

	return cmd
}
