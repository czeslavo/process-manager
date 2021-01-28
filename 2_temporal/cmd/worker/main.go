package main

import (
	"log"

	"github.com/czeslavo/process-manager/2_temporal/voiding"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, voiding.VoidDocumentsWorkflowQueueName, worker.Options{})

	w.RegisterWorkflowWithOptions(voiding.VoidDocuments, workflow.RegisterOptions{Name: voiding.VoidDocumentsWorkflowName})
	w.RegisterActivity(voiding.MarkDocumentAsVoidedInReports)
	w.RegisterActivity(voiding.VoidDocumentInBilling)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
