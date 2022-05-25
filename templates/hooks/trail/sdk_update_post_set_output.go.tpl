	if delta.DifferentAt("Spec.Tags") {
		err = rm.syncTrailTags(ctx, latest, desired)
		if err != nil {
			return nil, err
		}
	}