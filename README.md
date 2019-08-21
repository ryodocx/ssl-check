# ssl-check

## usage

```bash
# install
go get -u github.com/ryodocx/ssl-check

# usage
echo example.com | ssl-check              # output to stdout
echo example.com | ssl-check > result.csv # output to csv file

cat <<EOS > hostlist
> example.com
> example.com2
> example.com3
> EOS

cat hostlist | ssl-check > result.csv
ssl-check < hostlist > result.csv
```
