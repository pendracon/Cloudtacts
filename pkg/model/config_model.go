// Package model defines the structures and constants used by the application.
// This file contains query keys for accessing configuration options.
package model

const (
	KEY_CONFIG_FILE = "configFileId"

	KEY_CLIENT_COMMAND     = "commandId"
	KEY_CLIENT_IMAGE_FILE  = "imageId"
	KEY_CLIENT_IMAGE_TYPE  = "imageTypeId"
	KEY_CLIENT_INPUT_FILE  = "inputId"
	KEY_CLIENT_OUTPUT_FILE = "outputId"

	KEY_CLOUD_REGION  = "cloudRegionId"
	KEY_CLOUD_PROJECT = "cloudProjectId"

	KEY_AUTH_FUNCTION_HOST = "userdbFunctionHostId"
	KEY_AUTH_FUNCTION_PORT = "userdbFunctionPortId"
	KEY_AUTH_FUNCTION_GET  = "userdbGetUserId"
	KEY_AUTH_FUNCTION_ADD  = "userdbAddUserId"
	KEY_AUTH_FUNCTION_DEL  = "userdbDeleteUserId"
	KEY_AUTH_FUNCTION_UPD  = "userdbUpdateUserId"
	KEY_AUTH_FUNCTION_VAL  = "userdbValidateUserId"

	KEY_USERDB_TEST_MODE = "userdbTestModeId"
	KEY_USERDB_HOST_IP   = "userdbHostId"
	KEY_USERDB_PORT_NUM  = "userdbPortId"
	KEY_USERDB_DATABASE  = "userdbDatabaseId"
	KEY_USERDB_LOGIN     = "userdbLoginId"
	KEY_USERDB_PASSWORD  = "userdbCredsId"
	KEY_USERDB_MAX_POOL  = "userdbMaxPoolConnectionsId"
	KEY_USERDB_MAX_IDLE  = "userdbMaxIdleConnectionsId"
	KEY_USERDB_MAX_IDTM  = "userdbMaxIdleTimeId"
	KEY_USERDB_MAX_LFTM  = "userdbMaxLifeTimeId"
	KEY_STORAGE_BUCKET   = "storageBucketNameId"
)
