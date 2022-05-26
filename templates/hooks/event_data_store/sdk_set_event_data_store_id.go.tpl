    // This is a temporary solution.
    // For ReadOne and delete operations EventDataStore should be populated with
    // the ARN or the ID (which is not returned by the API).
    // Ideally we should be able to instruct the code generator to set an operation's
    // request field from a specific CR field path.
    input.EventDataStore = (*string)(r.ko.Status.ACKResourceMetadata.ARN)