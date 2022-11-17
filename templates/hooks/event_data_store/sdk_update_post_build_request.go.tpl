    // For Read, Update and Delete operations, EventDataStore should be populated with
    // the ARN or the ID (ID is not returned by the API).
    input.EventDataStore = (*string)(desired.ko.Status.ACKResourceMetadata.ARN)
    