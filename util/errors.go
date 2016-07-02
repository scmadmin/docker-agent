package util

func AppendError(errs []error, err error) []error {
	//todo: this is all somewhat questionable...
	if (errs == nil) {
		errs = make([]error, 0)
	}
	errs = append(errs, err)
	return errs
}
