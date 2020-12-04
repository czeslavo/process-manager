This repository is a place for experimenting with a concept of process manager concept (in CQRS context). 

# What is a process manager?



## Examples 

### [Voiding](1_voiding)
In this example I'm trying to implement documents voiding process which involves two remote services communication. 

* There is a service called **Billing** which is responsible for issuing various billing documents (invoices, receipts etc.). 
* There is a service called **Reports** which is responsible for aggregating documents issued by Billing in form 
of a report. A report is generated for every customer separately. 

When Billing issues a document, it emits a domain event **DocumentIssued**:

```json 
{
    "document_id": "document_unique_id",
    "client_id": "client_unique_id",
    "total_amount": "25.00"
}
```

Reports service listens to those events and builds a report based on them.  

Problem arises when someone decides that a particular document should be voided. A rule which cannot be broken
is that a document cannot be voided if a report it's in has already been sent to a customer. 

A user issues a **VoidDocument** command to Billing service which emits **DocumentVoidingRequested**   

#### Process prerequisites: 
1. **DocumentIssued** event from Billing 

#### Process flow:
1. **VoidDocument** command to Billing
2. **DocumentVoidingRequested** event from Billing   
3. **MarkDocumentAsVoided** command to Reports 
4. **MarkingDocumentAsVoidedSucceeded/Failed** event from Reports (in case report wasn't sent yet/was sent already)
5. **Abort/CompleteDocumentVoiding** command to Billing 
