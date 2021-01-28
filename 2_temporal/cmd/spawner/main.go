package main

import (
	"context"
	"fmt"
	"log"

	"github.com/czeslavo/process-manager/2_temporal/voiding"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("couldn't create temporal client", err)
	}
	defer c.Close()

	id := uuid.NewString()
	options := client.StartWorkflowOptions{
		ID:        id,
		TaskQueue: voiding.VoidDocumentsWorkflowQueueName,
	}

	documentUUID := "my-document-uuid"
	we, err := c.ExecuteWorkflow(context.Background(), options, voiding.VoidDocuments, documentUUID)
	if err != nil {
		log.Fatalln("couldn't execute void documents workflow", err)
	}

	// Don't use we.Get as it's long running.
	fmt.Printf("Spawened workflow ID=%s with RunID=%s\n", we.GetID(), we.GetRunID())
}
