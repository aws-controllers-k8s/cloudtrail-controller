    if r.ko.Status.ACKResourceMetadata == nil ||
        r.ko.Status.ACKResourceMetadata.ARN == nil {
        // If we don't return a ackerr.NotFound here, the API will return an
        // invalid ID or ARN error.
        return nil, ackerr.NotFound
        // Another solution would be to set the input id to "0-0-0-0-0"
        // Fun fact, the cloudtrail API generates UUID v4 for EventDataStores. When
        // trying to query an EDS using an empty ID the API returns an InvalidARN 
        // Exception. However trying to get an EDS using 0-0-0-0-0 doesn't return
        // a InvalidIDOrARN exception... Yes, it is considered as a valid UUID v4).
        // So we can use this impossible to generate UUID v4 to guaranty that we
        // will get a 404 error.
    }
    // For ReadOne and delete operations EventDataStore should be populated with
    // the ARN or the ID (ID is not returned by the API).
    input.EventDataStore = (*string)(r.ko.Status.ACKResourceMetadata.ARN)
    