package main

func isLocalBooking(isLocal bool) bool {
	return isLocal || global.ApiEndpoint == "" || global.RemoteBookingEnabled != true
}
