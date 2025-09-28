Ref: [RFC 2865](https://datatracker.ietf.org/doc/html/rfc2865)


in a simple radius PAP authentication:

1. AccessRequest

    contains:

    1. Code             - 1 byte - 01
    2. Identifier       - 1 byte - 0-255; random (usually incremental per request chain; not in PAP); echeod back for tracking
    3. Length           - 2 byte - End computed ( replace with 0x00\0x00)
    4. Authenticator    - 16 byte - if simple authenticator it is random generated;
                                    if message authenticator is computed at end using HMAC-MD5 and secret as key; typically not used in PAP
    5. AVPs:

        5.1 Type    - 1 byte - refer RFC
        5.2 Length  - 1 byte - computed at end
        5.3 Value   - Variable - encoded     


    Password encoding:

    Password is an AVP ( User-Password ); encoding format as per rfc
    
    b1 = MD5(S + RA)       c(1) = p1 xor b1
    b2 = MD5(S + c(1))     c(2) = p2 xor b2
            .                       .
            .                       .
            .                       .
    bi = MD5(S + c(i-1))   c(i) = pi xor bi

    The String will contain c(1)+c(2)+...+c(i) where + denotes concatenation.