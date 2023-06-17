# Receipt-Processor
Simple API for the challenge located at https://github.com/fetch-rewards/receipt-processor-challenge.
First time using Go so please don't judge too hashly.

## Endpoints
### process receipts
* Path: /receipts/process
* Method: POST
* Payload: Receipt JSON
* Response: JSON containing an id for the receipt.

### get points
* Path: /receipts/{id}/points
* Method: GET
* Response: A JSON object containing the number of points awarded.

### get receipts
* Path: /receipts
* Method: GET
* Response: JSON containing all receipts.
