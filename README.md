# grafctl

grafctl is a command line tool for interacting with Grafana. Its primarily intended to be used with
[foyle.io](http://foyle.io) and [Runme.dev](http://runme.dev).


## Quickstart

### Install

1. Download the latest release from the [releases page](https://github.com/jlewi/grafctl/releases)

### Create Base Resources

Define one or more base resources in your ~/.grafctl directory. These will be GrafanaLink resources that 
define a template for each dashboard you want to generate links for.

1. In Grafana, open the dashboard that you would like to generate links for
1. Configure the dashboard it serves as an example of the views you want to generate
1. In the UI click the share button to get a URL for the graph;
   * **Do not use the short link**
1. Use the `links parse` command to generate a base resource and save it to a file in your ~/.grafctl directory

   ```
   grafctl links parse --url=${URL} -o ~/.grafctl/${NAME}.yaml
   ```
      
   * By default the `GrafanaLink` resource is given the name `${NAME}` but you can override it 
     by specifying the `--name=${CUSTOMNAME}` flag

### Generate Links

Now you can generate links by defining a patch file that specifies the query and time range for the link. 

Here's an example 

```
cat <<EOF >/tmp/patch.yaml
template: somequery
query:
    builderOptions:
        logQuery: "service:app"
range: 
    from: "now-1h"
    to: "now"
EOF
grafctl links build -p /tmp/patch.yaml
```

This will print out a hyperlink to Grafana.

* **template** is the name of the GrafanaLink to use as the base resource
  * The name is stored in the yaml file
    ```yaml
    apiVersion: grafctl.foyle.io/v1alpha1
    kind: GrafanaLink
    metadata:
      name: sql
    ```
    
* **query** This is the patch to be applied to the query
   * The fields you need to set will be specific to your dashboard
* **range** specifies the time range for the query using grafana's
[syntax for relative time ranges](https://grafana.com/docs/grafana/latest/dashboards/use-dashboards/#time-units-and-relative-ranges).
By default relative time ranges are converted to absolute time ranges before constructing the URL in order to
generate a stable URL. If you want relative times in the url add the field `fixTime: false` to the patch.

grafctl will apply the patch to the first query in the template if your template contains more than one query.


