package fusionauth

import (
	"context"
	"encoding/json"

	"github.com/FusionAuth/go-client/pkg/fusionauth"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type SonyPSNConnectIdentityProviderBody struct {
	IdentityProvider fusionauth.SonyPSNIdentityProvider `json:"identityProvider"`
}

type SonyPSNAppConfig struct {
	ButtonText         string `json:"buttonText,omitempty"`
	ClientID           string `json:"client_id,omitempty"`
	ClientSecret       string `json:"client_secret,omitempty"`
	CreateRegistration bool   `json:"createRegistration"`
	Enabled            bool   `json:"enabled,omitempty"`
	Scope              string `json:"scope,omitempty"`
}

func resourceIDPSonyPSN() *schema.Resource {
	return &schema.Resource{
		CreateContext: createIDPSonyPSN,
		ReadContext:   readIDPSonyPSN,
		UpdateContext: updateIDPSonyPSN,
		DeleteContext: deleteIdentityProvider,
		Schema: map[string]*schema.Schema{
			"idp_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The ID to use for the new identity provider. If not specified a secure random UUID will be generated.",
				ValidateFunc: validation.IsUUID,
				ForceNew:     true,
			},
			"application_configuration": {
				Optional:    true,
				Type:        schema.TypeSet,
				Description: "The configuration for each Application that the identity provider is enabled for.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"button_text": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level button text.",
						},
						"client_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level client_id.",
						},
						"client_secret": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level client_secret.",
						},
						"create_registration": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Determines if a UserRegistration is created for the User automatically or not. If a user doesn’t exist in FusionAuth and logs in through an identity provider, this boolean controls whether or not FusionAuth creates a registration for the User in the Application they are logging into.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determines if this identity provider is enabled for the Application specified by the applicationId key.",
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "This is an optional Application specific override for the top level scope.",
						},
					},
				},
			},
			"button_text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level button text to use on the FusionAuth login page for this Identity Provider.",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level Sony PlayStation Network client id for your Application. This value is retrieved from the Sony PlayStation Network developer website when you setup your Sony PlayStation Network developer account.",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The top-level client secret to use with the Sony PlayStation Network Identity Provider when retrieving the long-lived token. This value is retrieved from the Sony PlayStation Network developer website when you setup your Sony PlayStation Network developer account.",
			},
			"debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if debug is enabled for this provider. When enabled, each time this provider is invoked to reconcile a login an Event Log will be created.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Determines if this provider is enabled. If it is false then it will be disabled globally.",
			},
			"lambda_reconcile_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The unique Id of the lambda to used during the user reconcile process to map custom claims from the external identity provider to the FusionAuth user.",
				ValidateFunc: validation.IsUUID,
			},
			"linking_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CreatePendingLink",
					"LinkAnonymously",
					"LinkByEmail",
					"LinkByEmailForExistingUser",
					"LinkByUsername",
					"LinkByUsernameForExistingUser",
					"Unsupported",
				}, false),
				Description: "The linking strategy to use when creating the link between the {idp_display_name} Identity Provider and the user.",
			},
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The top-level scope that you are requesting from Sony PlayStation Network.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createIDPSonyPSN(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPSonyPSN(data)

	b, err := json.Marshal(o)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := createIdentityProvider(b, client, data.Get("idp_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}

func readIDPSonyPSN(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := i.(Client)
	b, err := readIdentityProvider(data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	var ipb SonyPSNConnectIdentityProviderBody
	_ = json.Unmarshal(b, &ipb)

	return buildResourceDataFromIDPSonyPSN(data, ipb.IdentityProvider)
}

func updateIDPSonyPSN(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	o := buildIDPSonyPSN(data)

	b, err := json.Marshal(o)
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(Client)
	bb, err := updateIdentityProvider(b, data.Id(), client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = json.Unmarshal(bb, &o)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(o.IdentityProvider.Id)
	return nil
}

func buildIDPSonyPSN(data *schema.ResourceData) SonyPSNConnectIdentityProviderBody {
	o := fusionauth.SonyPSNIdentityProvider{
		BaseIdentityProvider: fusionauth.BaseIdentityProvider{
			Debug:      data.Get("debug").(bool),
			Enableable: buildEnableable("enabled", data),
			LambdaConfiguration: fusionauth.ProviderLambdaConfiguration{
				ReconcileId: data.Get("lambda_reconcile_id").(string),
			},
			Type:            fusionauth.IdentityProviderType_SonyPSN,
			LinkingStrategy: fusionauth.IdentityProviderLinkingStrategy(data.Get("linking_strategy").(string)),
		},
		ButtonText:   data.Get("button_text").(string),
		ClientId:     data.Get("client_id").(string),
		ClientSecret: data.Get("client_secret").(string),
		Scope:        data.Get("scope").(string),
	}

	ac := buildSonyPSNAppConfig("application_configuration", data)
	o.ApplicationConfiguration = ac
	return SonyPSNConnectIdentityProviderBody{IdentityProvider: o}
}

func buildResourceDataFromIDPSonyPSN(data *schema.ResourceData, res fusionauth.SonyPSNIdentityProvider) diag.Diagnostics {
	if err := data.Set("button_text", res.ButtonText); err != nil {
		return diag.Errorf("idpSonyPSN.button_text: %s", err.Error())
	}
	if err := data.Set("client_id", res.ClientId); err != nil {
		return diag.Errorf("idpSonyPSN.client_id: %s", err.Error())
	}
	if err := data.Set("client_secret", res.ClientSecret); err != nil {
		return diag.Errorf("idpSonyPSN.client_secret: %s", err.Error())
	}
	if err := data.Set("debug", res.Debug); err != nil {
		return diag.Errorf("idpSonyPSN.debug: %s", err.Error())
	}
	if err := data.Set("enabled", res.Enabled); err != nil {
		return diag.Errorf("idpSonyPSN.enabled: %s", err.Error())
	}
	if err := data.Set("lambda_reconcile_id", res.LambdaConfiguration.ReconcileId); err != nil {
		return diag.Errorf("idpSonyPSN.lambda_reconcile_id: %s", err.Error())
	}
	if err := data.Set("linking_strategy", res.LinkingStrategy); err != nil {
		return diag.Errorf("idpExternalJwt.linking_strategy: %s", err.Error())
	}
	if err := data.Set("scope", res.Scope); err != nil {
		return diag.Errorf("idpSonyPSN.scope: %s", err.Error())
	}

	// Since this is coming down as an interface and would end up being map[string]interface{}
	// with one of the values being map[string]interface{}
	b, _ := json.Marshal(res.ApplicationConfiguration)
	m := make(map[string]SonyPSNAppConfig)
	_ = json.Unmarshal(b, &m)

	ac := make([]map[string]interface{}, 0, len(res.ApplicationConfiguration))
	for k, v := range m {
		ac = append(ac, map[string]interface{}{
			"application_id":      k,
			"button_text":         v.ButtonText,
			"client_id":           v.ClientID,
			"client_secret":       v.ClientSecret,
			"create_registration": v.CreateRegistration,
			"enabled":             v.Enabled,
			"scope":               v.Scope,
		})
	}
	if err := data.Set("application_configuration", ac); err != nil {
		return diag.Errorf("idpSonyPSN.application_configuration: %s", err.Error())
	}

	return nil
}

func buildSonyPSNAppConfig(key string, data *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	s := data.Get(key)
	set, ok := s.(*schema.Set)
	if !ok {
		return m
	}
	l := set.List()
	for _, x := range l {
		ac := x.(map[string]interface{})
		aid := ac["application_id"].(string)
		oc := SonyPSNAppConfig{
			ButtonText:         ac["button_text"].(string),
			ClientID:           ac["client_id"].(string),
			ClientSecret:       ac["client_secret"].(string),
			CreateRegistration: ac["create_registration"].(bool),
			Enabled:            ac["enabled"].(bool),
			Scope:              ac["scope"].(string),
		}
		m[aid] = oc
	}
	return m
}
