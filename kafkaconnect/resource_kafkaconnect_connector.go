package kafkaconnect

import (
	"github.com/go-kafka/connect"
	"github.com/hashicorp/terraform/helper/schema"
)

func newConnector() *schema.Resource {
	return &schema.Resource{
		Create: createConnector,
		Read:   readConnector,
		Update: updateConnector,
		Delete: deleteConnector,
		Exists: checkIfConnectorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func createConnector(data *schema.ResourceData, context interface{}) error {
	connector := buildConnector(data)
	client := context.(*connect.Client)

	if _, err := client.CreateConnector(connector); err != nil {
		return err
	}

	data.SetId(connector.Name)

	return readConnector(data, context)
}

func readConnector(data *schema.ResourceData, context interface{}) error {
	client := context.(*connect.Client)
	connector, _, err := client.GetConnector(data.Id())
	if err != nil {
		return err
	}

	data.Set("name", connector.Name)

	configuration := make(map[string]string)

	if len(configuration) != 0 {
		data.Set("configuration", configuration)
	}

	return nil
}

func updateConnector(data *schema.ResourceData, context interface{}) error {
	connector := buildConnector(data)
	client := context.(*connect.Client)

	_, _, err := client.UpdateConnectorConfig(data.Id(), connector.Config)
	if err != nil {
		return err
	}

	return readConnector(data, context)
}

func deleteConnector(data *schema.ResourceData, context interface{}) error {
	client := context.(*connect.Client)

	_, err := client.DeleteConnector(data.Id())
	if err != nil {
		return err
	}

	return nil
}

func checkIfConnectorExists(
	data *schema.ResourceData,
	context interface{}) (bool, error) {
	client := context.(*connect.Client)

	_, _, err := client.GetConnectorStatus(data.Id())
	if err != nil {
		if apiError, ok := err.(connect.APIError); ok {
			if apiError.Code == 404 {
				return false, nil
			}
		}

		return false, err
	}

	return true, nil
}

func buildConnector(d *schema.ResourceData) *connect.Connector {
	connectorConfig := connect.ConnectorConfig{
		"name": d.Get("name").(string),
	}

	if configuration, ok := d.GetOk("configuration"); ok {
		for key, val := range configuration.(map[string]interface{}) {
			connectorConfig[key] = val.(string)
		}
	}

	connector := &connect.Connector{
		Name:   d.Get("name").(string),
		Config: connectorConfig,
	}

	return connector
}
