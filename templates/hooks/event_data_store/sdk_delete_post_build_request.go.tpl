    // For ReadOne and delete operations EventDataStore should be populated with
    // the ARN or the ID (ID is not returned by the API).
    input.EventDataStore = (*string)(r.ko.Status.ACKResourceMetadata.ARN)
    