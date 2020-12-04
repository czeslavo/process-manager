This repository is a place for experimenting with a concept of process manager (in CQRS context). 

# What is a process manager?

TODO

## Examples 

### [Voiding](1_voiding)
In this example I'm trying to implement documents voiding process which involves two remote services communication. 

* There is a service called **Billing** which is responsible for issuing various billing documents (invoices, receipts etc.). 
* There is a service called **Reports** which is responsible for aggregating documents issued by Billing in form 
of a report. A report is generated for every customer separately. 

When Billing issues a document, it emits a domain event **DocumentIssued**. Reports service listens to those events and builds a report based on them.  

Problem arises when someone decides that a particular document should be voided. A rule which cannot be broken
is that a document cannot be voided if a report in which the document is included has already been sent to a customer (published). 

A user issues a **VoidDocument** command to Billing service which emits **DocumentVoidingRequested**, what triggers the process.
There can be no more than 1 voiding process running for a single document at one time.



#### Process prerequisites
1. **DocumentIssued** event from Billing 

#### Process flow
1. **VoidDocument** command to Billing
2. **DocumentVoidingRequested** event from Billing   
3. **MarkDocumentAsVoided** command to Reports 
4. **MarkingDocumentAsVoidedSucceeded/Failed** event from Reports (in case report wasn't sent yet/was sent already)
5. **Abort/CompleteDocumentVoiding** command to Billing 
6. *(optional)* **AcknowledgeProcessFailure** command to process manager

A process has its identifier which is passed along with commands and events related to it (as a correlation id). 

What's important - the process manager creates process instances what makes them first class citizens - 
they're very similar to any different aggregate. They keep their own state and make the process 
visible and inspectable. In this particular example I made use of that and provided a command which 
allows to acknowledge a process failure (making it disappear from the view and allowing triggering new process for a particular document).

#### State machine describing the process
```
 +------------------+      +-------------------------------+          +----------------------+
 | VoidingRequested | +--> | MarkingDocumentAsVoidedFailed | +------> | Failure Acknowledged |
 +------------------+      +-------------------------------+          +----------------------+
         |
         +                 +----------------------------------+          +----------------+
         +---------------> | MarkingDocumentAsVoidedSucceeded | +------> | DocumentVoided |
                           +----------------------------------+          +----------------+
```


It's possible to play with the system using a simplistic HTML interface. To run a server: 
```
$ cd 1_voiding && go run .
``` 

In order to make the process last longer, so you could observe the progress, run the server as following:
```
$ SLOW_DOWN=1 cd 1_voiding && go run .
``` 


The interface will be available under `http://localhost:8080`.

