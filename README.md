# The Bridge 

Development Dependencies:
```
brew install golang-migrate
brew install make

```

# Todos :

[] Add a seperate api for image upload, which will be called by fronted and then it will take the res (image url) and send to actual api


# TEMP :

```

-- Check if A follows B (O(1) PK lookup)
SELECT 1 FROM user_follows WHERE follower_id = $1 AND followee_id = $2;

-- Get all followers of a user
SELECT follower_id FROM user_follows WHERE followee_id = $1;

-- Mutual follows (people you follow who also follow you back)
SELECT follower_id FROM user_follows
WHERE followee_id = $me
  AND follower_id IN (SELECT followee_id FROM user_follows WHERE follower_id = $me);

```




## Developer Resources

- [Golang-Migrate](https://dev.to/wiliamvj/using-migrations-with-golang-3449)
- [AUTH] (https://web.archive.org/web/20240917031521/https://hrishikeshpathak.com/blog/complete-guide-to-oauth2-in-golang/)
- 


# Notes :

1. handle this auth flow:

Unless frontend does this:

```
if (response.headers["x-new-access-token"]) {
    localStorage.setItem("access_token", response.headers["x-new-access-token"]);
}
```

Then:

Next request:
```
Authorization: Bearer <new_access_token>
```

If frontend ignores header?

Then next request will again use old expired token → middleware refreshes again → loop forever.


# TODOs (For V2) : 

1. Improve the security by rotating refresh tokens eveytime we generate a new access token. This way, if a refresh token is compromised, it can only be used once before it becomes invalid. and also implement a mechanism to revoke refresh tokens when necessary, such as when a user logs out or when suspicious activity is detected. This can be done by maintaining a blacklist of revoked tokens or by using a token versioning system where each token has a version number that is incremented upon rotation.

Add : 
```
Session table

Device tracking

Token hashing

Rotation strategy

Logout-all-devices

Session revocation logic
```

right now we don't have these features, but they are important for a secure and robust authentication system. Implementing these features will help protect user accounts and prevent unauthorized access in case of token compromise.

1. Per-device logout
2. Logout-all-devices
3. Immediate refresh revocation
4. Stolen refresh detection