# Changelog

## 0.6.2

- fix: Copy traces only if source is errx.Error type

## 0.6.1

- fix: Nil error tracing

## 0.6.0

- feat(error): Add option to AddMetadata when tracing

## 0.5.1

- fix(error): Merge traces on wrapping errx.Error
- fix(error): Write CausedBy after Traces

## 0.5.0

- feat(error): Add instance method AddMetadata on errx.Error

## 0.4.0

- feat(error): Add Message getter function

## 0.3.0

- feat(option): Add option to trace generic error with fmt.Errorf

## 0.2.1

- fix(error): Rename TraceError to Trace; Make source as optional parameter

## 0.2.0

- feat(option): Add AddMetadata error option function
- feat(builder): Add CopyError function to Copy errx.Error and override namespace

## 0.1.0

- feat(builder): Add error builder
- feat(error): Add errx.Error
