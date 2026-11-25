
## OIDC Single Sign On

AELIG Mail supports single sign-on with OIDC (OpenID Connect). Any standards compliant OIDC provider can be configured in Settings -> Security -> OIDC

Enable **OIDC RP-initiated logout** in Settings → Security to redirect to the provider's logout endpoint when logging out of an OIDC session, if the provider supports OIDC logout. This option is disabled by default. Register `https://mail.yoursite.com/admin/login` (including any configured root URL path) as an allowed post-logout redirect URI with the provider.

### User auto-creation
If `Settings -> Security -> OIDC -> Auto-create users` is turned on, when users login via OIDC, an account is auto-created if an existing account is not found (based on the OIDC e-mail ID). If role mapping is enabled, newly created users (on first OIDC login) can be auto-assigned roles mapped to OIDC claims.

# Tutorials

Tutorials for configuring AELIG Mail SSO with popular OIDC providers.

## Keycloak
Keycloak configuration for AELIG Mail SSO integration.

### 1. Create a new client in Keycloak
In the Keycloak admin, use an existing realm, or create a new realm. Create a new client in `Clients → Create`.

- **General Settings**
    - **Client type**: `OpenID Connect`
    - **Client ID**: `aelig-mail` (or any preferred name)
    - **Name**: Optional descriptive name (e.g., "AELIG Mail SSO")
- **Capability Config**:
    - **Client authentication**: On
    - **Authorization**: On
    - **Authentication Flow**
        - **Standard Flow**: On
        - **Direct Access grants**: On
- **Login Settings**:
    - **Root URL**: Copy the **Redirect URL for oAuth provider** value from AELIG Mail Admin -> Settings -> Security -> OIDC. It will look like `https://mail.yoursite.com/auth/oidc`
    - **Valid redirect URIs**: Same as the Root URL above
    - **Valid post logout redirect URIs**: *

After the client creation steps above, go to the client's `Credentials` tab and copy the `Client Secret`.

### 2. Configure AELIG Mail
2. In AELIG Mail Admin -> Settings -> Security -> OIDC.
    - **Enable OIDC SSO**: Turn on
    - **Provider URL**: `https://keycloak.yoursite.com/realms/{realm}` (replace `{realm}` with the chosen realm name). This URL is as of v26.3 and may differ across Keycloak versions.
    - **Provider name**: Set a name to show on the AELIG Mail login form, eg: `Login with OrgName`
    - **Client ID**: Client ID set in Keycloak, eg: `aelig-mail`
    - **Client Secret**: Client Secret copied from Keycloak
    - **Auto-create users from SSO**: (Optional) Enable to automatically create users who don't exist
    - **Default user role**: (Required if auto-create enabled) Select role for new users



## Authentik
Authentik configuration for AELIG Mail SSO integration.

### 1. Create a new OIDC provider in Authentik
In the Authentik admin interface, create a new OIDC provider for AELIG Mail.

- **Provider Settings**:
    - **Name**: `aelig-mail` (or any preferred name)
    - **Signing Key**: `authentik Self-signed Certificate`
    - **Client Type**: `Confidential`
    - **Client ID**: `aelig-mail` (or any preferred name)
    - **Redirect URIs**: Copy the **Redirect URL for oAuth provider** value from AELIG Mail Admin -> Settings -> Security -> OIDC. It will look like `https://mail.yoursite.com/auth/oidc`

After creating the provider, copy the **Client Secret**.

### 2. Create an application in Authentik
Create a new application and connect it to the newly created provider.

- **Application Settings**:
    - **Name**: `aelig-mail` (or any preferred name)
    - **Slug**: `aelig-mail` (or any preferred slug. Used in the redirect URL)
    - **Provider**: Select the OIDC provider created in the previous step

### 3. Configure AELIG Mail
In AELIG Mail Admin → Settings → Security → OIDC:

- **Enable OIDC SSO**: Turn on
- **Provider URL**: `https://authentik.yoursite.com/application/o/{slug}/` (replace `{slug}` with the application's slug)
- **Provider Name**: Set a name to show on the login form (e.g., `Login with OrgName`)
- **Client ID**: Client ID set in Authentik (e.g., `aelig-mail`)
- **Client Secret**: Client Secret copied from Authentik
- **Auto-create users from SSO**: (Optional) Enable to automatically create users who don't exist
- **Default user role**: (Required if auto-create enabled) Select role for new users

## Google Workspace
Google Workspace (Google Cloud) configuration for AELIG Mail SSO integration.

### 1. Create a new OIDC provider in Google Cloud Console / Google Workspace
In the Google Cloud Console interface, create a new Project.

- **Project Settings**:
    - **Project name**: `aelig-mail` (or any preferred name)
- **Branding Settings**:
    - **App name**: `aelig-mail` (or any preferred name, this will be visible to the users.)
    - **Authorised domains**: `mail.example.com` (or domains that your instance is available on.)

After creating the project, goto **Clients**.

### 2. Create an client in project.
Create a new client and configure it.

- **Application Settings**:
    - **Application type**: `Web application`
    - **Name**: `aelig-mail` (or any preferred name)
    - **Authorised JavaScript origins**: `https://mail.example.com` (or domains that your instance is available on.)
    - **Authorised redirect URIs**: `https://mail.example.com/auth/oidc` (or domains that your instance is available on, value is also available in the Settings mentioned above. (Redirect URL for oAuth provider))

Hit save and note the Client ID and Client Secret

### 3. Configure AELIG Mail
In AELIG Mail Admin → Settings → Security → OIDC:

- **Enable OIDC SSO**: Turn on
- **Provider URL**: `https://accounts.google.com` (select Google to Auto-Fill)
- **Provider Name**: Set a name to show on the login form (e.g., `Login with OrgName`)
- **Client ID**: Client ID copied from Console (e.g., `XXXX.apps.googleusercontent.com`)
- **Client Secret**: Client Secret copied from Console
- **Auto-create users from SSO**: (Optional) Enable to automatically create users who don't exist
- **Default user role**: (Required if auto-create enabled) Select role for new users
