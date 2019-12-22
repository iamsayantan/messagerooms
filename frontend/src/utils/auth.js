var accessTokenKey = "accessToken"
var authUserKey = "authenticatedUser"

export function storeAccessToken(token) {
    localStorage.setItem(accessTokenKey, token)
}

export function getAccessToken() {
    return localStorage.getItem(accessTokenKey)
}

export function storeUser(user) {
    localStorage.setItem(authUserKey, JSON.stringify(user))
}

export function getAuthenticatedUser() {
    return JSON.parse(localStorage.getItem(authUserKey))
}