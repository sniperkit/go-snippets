// Default configuration settings.  Values will be overriden
// at runtime if they are defined locally by the user.
//
// See options.html and options.js for more information
//
sdes.config = {
    "routingRules" : [ 
        {
            "matches": "https://github.com/*",
            "host": {
                "api": "https://api.github.com",
                "baseUrl": "https://github.com",
                "secret": "",
                "type": "github",
                "xFrameOptions": "DENY"
            },
            "gitsense": {
                "baseUrl": "https://public.gitsense.com",
                "secret": ""
            }
        },
        {
            "matches": "https://gitlab.com/*",
            "host": {
                "api": "https://gitlab.com/api/v3",
                "baseUrl": "https://gitlab.com",
                "secret": "",
                "type": "gitlab",
                "xFrameOptions": "DENY"
            },
            "gitsense": {
                "baseUrl": "https://public.gitsense.com",
                "secret": ""
            }
        }
    ]
}
