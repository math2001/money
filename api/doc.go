/*
Data folder structure:

    appsalts
    users/
        {id1}/
            kmdata/                # this is managed by keys manager
                keys
                passwordhashfile   # in the user's case, this file is redundant
                salts
            recordfile1            # these are encrypted with kmdata just above
            recordfile2
            ...
        {id2}/
            kmdata/
                keys
                passwordhashfile
                salts
            recordfile1
            recordfile2
            ...


    The appsalts contains two salts:

    saltcookie                 # the salt used to generate the keys used to
                               # sign the cookies
    saltpassword               # the salt used to encrypt the passwords
                               # within the database
*/
package api
