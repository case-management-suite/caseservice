package caseservice

// type EngineQueueConsumerHandler struct {
// 	casesService CaseService
// }

// func NewEngineQueueConsumerHandler(casesService *CaseService) queue.QueueHandler {
// 	return EngineQueueConsumerHandler{casesService: *casesService}
// }

// type ParsedQueueEvent struct {
// 	UUID       models.Identifier
// 	Job        string
// 	CaseAction models.CaseAction
// }

// func ParseQueueEvent(job []byte) (ParsedQueueEvent, error) {
// 	record := models.CaseAction{}
// 	err := json.Unmarshal(job, &record)
// 	return ParsedQueueEvent{CaseAction: record, UUID: job.UUID, Job: job.Job}, err
// }

// func (cc EngineQueueConsumerHandler) ConsumeEvents(event []byte) error {
// 	d, err := ParseQueueEvent(event)
// 	if err != nil {
// 		log.Err(err).Str("function", "queue_handler.ConsumeCaseEvents").Msg("Could not parse the queue event")
// 		return err
// 	}
// 	log.Printf("Received a message: %s", d.Job)

// 	cc.casesService.UpdateCase(d.CaseAction.CaseRecord)
// 	return nil
// }
