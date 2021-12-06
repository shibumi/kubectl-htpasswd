# kubectl-htpasswd

kubectl-htpasswd is a nginx-ingress compatible basic-auth secret generator.
Although Kubernetes supports basic-auth secrets, these secrets are incompatible with
many ingresses such like NGINX ingress. This plugin provides an easy way to create
bcrypt hashed secrets on the fly without much hassle.

### Supported hash algorithms

* bcrypt

## Examples

### Create the secret on the cluster in the current namespace

`$ kubectl htpasswd create $SECRETNAME $USER1=$PASSWORD1 $USER2=$PASSWORD2`

### Just print the secret in yaml

`$ kubectl htpasswd create $SECRETNAME $USER1=$PASSWORD1 $USER2=$PASSWORD2 -o yaml --dry-run`
`

