package audit_manager

import (
	"context"

	"github.com/Velocidex/ordereddict"
	"github.com/sirupsen/logrus"
	config_proto "www.velocidex.com/golang/velociraptor/config/proto"
	"www.velocidex.com/golang/velociraptor/logging"
	"www.velocidex.com/golang/velociraptor/services"
)

type AuditManager struct{}

func (self *AuditManager) LogAudit(
	ctx context.Context,
	config_obj *config_proto.Config,
	principal, operation string,
	details *ordereddict.Dict) error {

	record := ordereddict.NewDict().
		Set("operation", operation).
		Set("principal", principal).
		Set("details", details)

	logger := logging.GetLogger(config_obj, &logging.Audit)
	logger.WithFields(logrus.Fields(*record.ToDict())).Info(operation)

	journal, err := services.GetJournal(config_obj)
	if err != nil {
		return err
	}

	// If an event is important enough to be audit logged we need to
	// make sure to write it syncronously.
	return journal.PushRowsToArtifact(
		ctx, config_obj, []*ordereddict.Dict{record},
		"Server.Audit.Logs", "server", "")
}
