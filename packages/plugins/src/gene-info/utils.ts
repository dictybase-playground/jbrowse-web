import { pipe } from "fp-ts/function"
import { compose } from "fp-ts/Refinement"

const isNonNullableObject = (x: unknown): x is NonNullable<object> =>
  typeof x === "object" && x !== null
const hasTrackType = (obj: NonNullable<object>): obj is { trackType: string } =>
  "trackType" in obj
const hasType = (obj: NonNullable<object>): obj is { type: string } =>
  "type" in obj

const unknownHasTrackType = pipe(isNonNullableObject, compose(hasTrackType))
const unknownHasType = pipe(isNonNullableObject, compose(hasType))

export { unknownHasTrackType, unknownHasType }
