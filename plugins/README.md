# Utils for agilestack plugins
This package contains the tools needed for the agilestack plugins to :

| Operation                                                           | Status                                                                   |
|---------------------------------------------------------------------|--------------------------------------------------------------------------|
| register themselves                                                 | TODO, check if needed, because registrator is already doing registration |
| register their resources to be secured, with associated permissions | TODO                                                                     |
| register the translations of the resources, permissions             | TODO                                                                     |

### Store in consul resources and permissions for security concerns :

The resources and their permissions are stored in consul 's KV store under : `/agilestack/security/resources/<resource>`

```
{
    "keys": {
        "lang": "accounts",
        "security": "rn:hydra:accounts"
    },
    "permissions": [
        "create",
        "get",
        "delete",
        "put:password",
        "put:data"
    ]
}
```

### Store in consul security translations for resources and permissions on security context :

The resources and permissions translations are stored in consul 's KV store under : `agileStack/translations/<lang>`

Each plugin comes with its own resources, they should be added to the global translation

**The global translation looks like this below (example for french)**
```
{
    "security": {
        "permissions: {
            "account": {
                "create": "Créer ou modifier"
            },
            "create": "Créer",
            "modify": "Modifier",
        },
        "resources": {
            "accounts": "comptes utilisateurs"
        },
        "create": "Créééééeer",
        "modify": "Modiiiiiiiiiiiifier"
    }
}
```
In this example :
 
* security, security.permissions, security.permissions.account are context of the translations
* create, modify, account are the keys to be translated, inside their context

This json format is compatible with tools such as https://poeditor.com