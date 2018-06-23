// Package leadership provides leadership election and observation.
//
// To participate in leadership election as a candidate, an application must
// get a candidate implementation from a LeadershipProvider. To observe the
// quorum, an application must instantiate get an Observer implementation from
// a LeadershipProvider.
package leadership
