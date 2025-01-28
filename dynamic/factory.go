package dynConfigs

import "time"

type FieldFactory interface {
	NewBoolean() BooleanEx
	NewInteger() IntegerEx
	NewFloat() FloatEx
	NewString() StringEx
	NewStrings() StringsEx
	NewDuration() DurationEx
	NewTime() TimeEx
}

type fieldFactory struct {
}

func (that *fieldFactory) NewBoolean() BooleanEx {
	return NewBoolean(false)
}

func (that *fieldFactory) NewInteger() IntegerEx {
	return NewInteger(0)
}

func (that *fieldFactory) NewFloat() FloatEx {
	return NewFloat(0)
}

func (that *fieldFactory) NewString() StringEx {
	return NewString("")
}

func (that *fieldFactory) NewStrings() StringsEx {
	return NewStrings([]string{})
}

func (that *fieldFactory) NewDuration() DurationEx {
	return NewDuration(0)
}

func (that *fieldFactory) NewTime() TimeEx {
	return NewTime(time.Time{})
}

func NewFieldFactory() FieldFactory {
	return &fieldFactory{}
}
