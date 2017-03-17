/*
Copyright 2016 The Archon Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloudinit

import (
	"github.com/coreos/yaml"
)

type CoreOSCloudConfig struct {
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	CoreOS            *CoreOS  `yaml:"coreos,omitempty"`
	WriteFiles        []File   `yaml:"write_files,omitempty"`
	Hostname          string   `yaml:"hostname,omitempty"`
	Users             []User   `yaml:"users,omitempty"`
	ManageEtcHosts    EtcHosts `yaml:"manage_etc_hosts,omitempty"`
}

func (cc CoreOSCloudConfig) Bytes() ([]byte, error) {
	data, err := yaml.Marshal(cc)
	if err != nil {
		return nil, err
	}
	return append([]byte("#cloud-config\n"), data...), nil
}

func (cc CoreOSCloudConfig) String() (string, error) {
	data, err := cc.Bytes()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type CoreOS struct {
	Etcd      *Etcd      `yaml:"etcd,omitempty"`
	Etcd2     *Etcd2     `yaml:"etcd2,omitempty"`
	Flannel   *Flannel   `yaml:"flannel,omitempty"`
	Fleet     *Fleet     `yaml:"fleet,omitempty"`
	Locksmith *Locksmith `yaml:"locksmith,omitempty"`
	OEM       *OEM       `yaml:"oem,omitempty"`
	Update    *Update    `yaml:"update,omitempty"`
	Units     []Unit     `yaml:"units,omitempty"`
}

type File struct {
	Encoding           string `yaml:"encoding,omitempty" valid:"^(base64|b64|gz|gzip|gz\\+base64|gzip\\+base64|gz\\+b64|gzip\\+b64)$"`
	Content            string `yaml:"content,omitempty"`
	Owner              string `yaml:"owner,omitempty"`
	Path               string `yaml:"path,omitempty"`
	RawFilePermissions string `yaml:"permissions,omitempty" valid:"^0?[0-7]{3,4}$"`
}

type User struct {
	Name                 string   `yaml:"name,omitempty"`
	PasswordHash         string   `yaml:"passwd,omitempty"`
	Sudo                 string   `yaml:"sudo,omitempty"`
	SSHAuthorizedKeys    []string `yaml:"ssh_authorized_keys,omitempty"`
	SSHImportGithubUser  string   `yaml:"coreos_ssh_import_github,omitempty"       deprecated:"trying to fetch from a remote endpoint introduces too many intermittent errors"`
	SSHImportGithubUsers []string `yaml:"coreos_ssh_import_github_users,omitempty" deprecated:"trying to fetch from a remote endpoint introduces too many intermittent errors"`
	SSHImportURL         string   `yaml:"coreos_ssh_import_url,omitempty"          deprecated:"trying to fetch from a remote endpoint introduces too many intermittent errors"`
	GECOS                string   `yaml:"gecos,omitempty"`
	Homedir              string   `yaml:"homedir,omitempty"`
	NoCreateHome         bool     `yaml:"no_create_home,omitempty"`
	PrimaryGroup         string   `yaml:"primary_group,omitempty"`
	Groups               []string `yaml:"groups,omitempty"`
	NoUserGroup          bool     `yaml:"no_user_group,omitempty"`
	System               bool     `yaml:"system,omitempty"`
	NoLogInit            bool     `yaml:"no_log_init,omitempty"`
	Shell                string   `yaml:"shell,omitempty"`
}

type EtcHosts string

type Etcd struct {
	Addr                     string  `yaml:"addr,omitempty"                          env:"ETCD_ADDR"`
	AdvertiseClientURLs      string  `yaml:"advertise_client_urls,omitempty"         env:"ETCD_ADVERTISE_CLIENT_URLS"       deprecated:"etcd2 options no longer work for etcd"`
	BindAddr                 string  `yaml:"bind_addr,omitempty"                     env:"ETCD_BIND_ADDR"`
	CAFile                   string  `yaml:"ca_file,omitempty"                       env:"ETCD_CA_FILE"`
	CertFile                 string  `yaml:"cert_file,omitempty"                     env:"ETCD_CERT_FILE"`
	ClusterActiveSize        int     `yaml:"cluster_active_size,omitempty"           env:"ETCD_CLUSTER_ACTIVE_SIZE"`
	ClusterRemoveDelay       float64 `yaml:"cluster_remove_delay,omitempty"          env:"ETCD_CLUSTER_REMOVE_DELAY"`
	ClusterSyncInterval      float64 `yaml:"cluster_sync_interval,omitempty"         env:"ETCD_CLUSTER_SYNC_INTERVAL"`
	CorsOrigins              string  `yaml:"cors,omitempty"                          env:"ETCD_CORS"`
	DataDir                  string  `yaml:"data_dir,omitempty"                      env:"ETCD_DATA_DIR"`
	Discovery                string  `yaml:"discovery,omitempty"                     env:"ETCD_DISCOVERY"`
	DiscoveryFallback        string  `yaml:"discovery_fallback,omitempty"            env:"ETCD_DISCOVERY_FALLBACK"          deprecated:"etcd2 options no longer work for etcd"`
	DiscoverySRV             string  `yaml:"discovery_srv,omitempty"                 env:"ETCD_DISCOVERY_SRV"               deprecated:"etcd2 options no longer work for etcd"`
	DiscoveryProxy           string  `yaml:"discovery_proxy,omitempty"               env:"ETCD_DISCOVERY_PROXY"             deprecated:"etcd2 options no longer work for etcd"`
	ElectionTimeout          int     `yaml:"election_timeout,omitempty"              env:"ETCD_ELECTION_TIMEOUT"            deprecated:"etcd2 options no longer work for etcd"`
	ForceNewCluster          bool    `yaml:"force_new_cluster,omitempty"             env:"ETCD_FORCE_NEW_CLUSTER"           deprecated:"etcd2 options no longer work for etcd"`
	GraphiteHost             string  `yaml:"graphite_host,omitempty"                 env:"ETCD_GRAPHITE_HOST"`
	HeartbeatInterval        int     `yaml:"heartbeat_interval,omitempty"            env:"ETCD_HEARTBEAT_INTERVAL"          deprecated:"etcd2 options no longer work for etcd"`
	HTTPReadTimeout          float64 `yaml:"http_read_timeout,omitempty"             env:"ETCD_HTTP_READ_TIMEOUT"`
	HTTPWriteTimeout         float64 `yaml:"http_write_timeout"            env:"ETCD_HTTP_WRITE_TIMEOUT"`
	InitialAdvertisePeerURLs string  `yaml:"initial_advertise_peer_urls"   env:"ETCD_INITIAL_ADVERTISE_PEER_URLS" deprecated:"etcd2 options no longer work for etcd"`
	InitialCluster           string  `yaml:"initial_cluster"               env:"ETCD_INITIAL_CLUSTER"             deprecated:"etcd2 options no longer work for etcd"`
	InitialClusterState      string  `yaml:"initial_cluster_state"         env:"ETCD_INITIAL_CLUSTER_STATE"       deprecated:"etcd2 options no longer work for etcd"`
	InitialClusterToken      string  `yaml:"initial_cluster_token"         env:"ETCD_INITIAL_CLUSTER_TOKEN"       deprecated:"etcd2 options no longer work for etcd"`
	KeyFile                  string  `yaml:"key_file"                      env:"ETCD_KEY_FILE"`
	ListenClientURLs         string  `yaml:"listen_client_urls"            env:"ETCD_LISTEN_CLIENT_URLS"          deprecated:"etcd2 options no longer work for etcd"`
	ListenPeerURLs           string  `yaml:"listen_peer_urls"              env:"ETCD_LISTEN_PEER_URLS"            deprecated:"etcd2 options no longer work for etcd"`
	MaxResultBuffer          int     `yaml:"max_result_buffer"             env:"ETCD_MAX_RESULT_BUFFER"`
	MaxRetryAttempts         int     `yaml:"max_retry_attempts"            env:"ETCD_MAX_RETRY_ATTEMPTS"`
	MaxSnapshots             int     `yaml:"max_snapshots"                 env:"ETCD_MAX_SNAPSHOTS"               deprecated:"etcd2 options no longer work for etcd"`
	MaxWALs                  int     `yaml:"max_wals"                      env:"ETCD_MAX_WALS"                    deprecated:"etcd2 options no longer work for etcd"`
	Name                     string  `yaml:"name"                          env:"ETCD_NAME"`
	PeerAddr                 string  `yaml:"peer_addr"                     env:"ETCD_PEER_ADDR"`
	PeerBindAddr             string  `yaml:"peer_bind_addr"                env:"ETCD_PEER_BIND_ADDR"`
	PeerCAFile               string  `yaml:"peer_ca_file"                  env:"ETCD_PEER_CA_FILE"`
	PeerCertFile             string  `yaml:"peer_cert_file"                env:"ETCD_PEER_CERT_FILE"`
	PeerElectionTimeout      int     `yaml:"peer_election_timeout"         env:"ETCD_PEER_ELECTION_TIMEOUT"`
	PeerHeartbeatInterval    int     `yaml:"peer_heartbeat_interval"       env:"ETCD_PEER_HEARTBEAT_INTERVAL"`
	PeerKeyFile              string  `yaml:"peer_key_file"                 env:"ETCD_PEER_KEY_FILE"`
	Peers                    string  `yaml:"peers"                         env:"ETCD_PEERS"`
	PeersFile                string  `yaml:"peers_file"                    env:"ETCD_PEERS_FILE"`
	Proxy                    string  `yaml:"proxy"                         env:"ETCD_PROXY"                       deprecated:"etcd2 options no longer work for etcd"`
	RetryInterval            float64 `yaml:"retry_interval"                env:"ETCD_RETRY_INTERVAL"`
	Snapshot                 bool    `yaml:"snapshot"                      env:"ETCD_SNAPSHOT"`
	SnapshotCount            int     `yaml:"snapshot_count"                env:"ETCD_SNAPSHOTCOUNT"`
	StrTrace                 string  `yaml:"trace"                         env:"ETCD_TRACE"`
	Verbose                  bool    `yaml:"verbose"                       env:"ETCD_VERBOSE"`
	VeryVerbose              bool    `yaml:"very_verbose"                  env:"ETCD_VERY_VERBOSE"`
	VeryVeryVerbose          bool    `yaml:"very_very_verbose"             env:"ETCD_VERY_VERY_VERBOSE"`
}

type Etcd2 struct {
	AdvertiseClientURLs      string `yaml:"advertise_client_urls,omitempty"         env:"ETCD_ADVERTISE_CLIENT_URLS"`
	CAFile                   string `yaml:"ca_file,omitempty"                       env:"ETCD_CA_FILE"                     deprecated:"ca_file obsoleted by trusted_ca_file and client_cert_auth"`
	CertFile                 string `yaml:"cert_file,omitempty"                     env:"ETCD_CERT_FILE"`
	ClientCertAuth           bool   `yaml:"client_cert_auth,omitempty"              env:"ETCD_CLIENT_CERT_AUTH"`
	CorsOrigins              string `yaml:"cors,omitempty"                          env:"ETCD_CORS"`
	DataDir                  string `yaml:"data_dir,omitempty"                      env:"ETCD_DATA_DIR"`
	Debug                    bool   `yaml:"debug,omitempty"                         env:"ETCD_DEBUG"`
	Discovery                string `yaml:"discovery,omitempty"                     env:"ETCD_DISCOVERY"`
	DiscoveryFallback        string `yaml:"discovery_fallback,omitempty"            env:"ETCD_DISCOVERY_FALLBACK"`
	DiscoverySRV             string `yaml:"discovery_srv,omitempty"                 env:"ETCD_DISCOVERY_SRV"`
	DiscoveryProxy           string `yaml:"discovery_proxy,omitempty"               env:"ETCD_DISCOVERY_PROXY"`
	ElectionTimeout          int    `yaml:"election_timeout,omitempty"              env:"ETCD_ELECTION_TIMEOUT"`
	ForceNewCluster          bool   `yaml:"force_new_cluster,omitempty"             env:"ETCD_FORCE_NEW_CLUSTER"`
	HeartbeatInterval        int    `yaml:"heartbeat_interval,omitempty"            env:"ETCD_HEARTBEAT_INTERVAL"`
	InitialAdvertisePeerURLs string `yaml:"initial_advertise_peer_urls,omitempty"   env:"ETCD_INITIAL_ADVERTISE_PEER_URLS"`
	InitialCluster           string `yaml:"initial_cluster,omitempty"               env:"ETCD_INITIAL_CLUSTER"`
	InitialClusterState      string `yaml:"initial_cluster_state,omitempty"         env:"ETCD_INITIAL_CLUSTER_STATE"`
	InitialClusterToken      string `yaml:"initial_cluster_token,omitempty"         env:"ETCD_INITIAL_CLUSTER_TOKEN"`
	KeyFile                  string `yaml:"key_file,omitempty"                      env:"ETCD_KEY_FILE"`
	ListenClientURLs         string `yaml:"listen_client_urls,omitempty"            env:"ETCD_LISTEN_CLIENT_URLS"`
	ListenPeerURLs           string `yaml:"listen_peer_urls,omitempty"              env:"ETCD_LISTEN_PEER_URLS"`
	LogPackageLevels         string `yaml:"log_package_levels,omitempty"            env:"ETCD_LOG_PACKAGE_LEVELS"`
	MaxSnapshots             int    `yaml:"max_snapshots,omitempty"                 env:"ETCD_MAX_SNAPSHOTS"`
	MaxWALs                  int    `yaml:"max_wals,omitempty"                      env:"ETCD_MAX_WALS"`
	Name                     string `yaml:"name,omitempty"                          env:"ETCD_NAME"`
	PeerCAFile               string `yaml:"peer_ca_file,omitempty"                  env:"ETCD_PEER_CA_FILE"                deprecated:"peer_ca_file obsoleted peer_trusted_ca_file and peer_client_cert_auth"`
	PeerCertFile             string `yaml:"peer_cert_file,omitempty"                env:"ETCD_PEER_CERT_FILE"`
	PeerKeyFile              string `yaml:"peer_key_file,omitempty"                 env:"ETCD_PEER_KEY_FILE"`
	PeerClientCertAuth       bool   `yaml:"peer_client_cert_auth,omitempty"         env:"ETCD_PEER_CLIENT_CERT_AUTH"`
	PeerTrustedCAFile        string `yaml:"peer_trusted_ca_file,omitempty"          env:"ETCD_PEER_TRUSTED_CA_FILE"`
	Proxy                    string `yaml:"proxy,omitempty"                         env:"ETCD_PROXY"                       valid:"^(on|off|readonly)$"`
	ProxyDialTimeout         int    `yaml:"proxy_dial_timeout,omitempty"            env:"ETCD_PROXY_DIAL_TIMEOUT"`
	ProxyFailureWait         int    `yaml:"proxy_failure_wait,omitempty"            env:"ETCD_PROXY_FAILURE_WAIT"`
	ProxyReadTimeout         int    `yaml:"proxy_read_timeout,omitempty"            env:"ETCD_PROXY_READ_TIMEOUT"`
	ProxyRefreshInterval     int    `yaml:"proxy_refresh_interval,omitempty"        env:"ETCD_PROXY_REFRESH_INTERVAL"`
	ProxyWriteTimeout        int    `yaml:"proxy_write_timeout,omitempty"           env:"ETCD_PROXY_WRITE_TIMEOUT"`
	SnapshotCount            int    `yaml:"snapshot_count,omitempty"                env:"ETCD_SNAPSHOT_COUNT"`
	TrustedCAFile            string `yaml:"trusted_ca_file,omitempty"               env:"ETCD_TRUSTED_CA_FILE"`
	WalDir                   string `yaml:"wal_dir,omitempty"                       env:"ETCD_WAL_DIR"`
}

type Flannel struct {
	EtcdEndpoints string `yaml:"etcd_endpoints,omitempty" env:"FLANNELD_ETCD_ENDPOINTS"`
	EtcdCAFile    string `yaml:"etcd_cafile,omitempty"    env:"FLANNELD_ETCD_CAFILE"`
	EtcdCertFile  string `yaml:"etcd_certfile,omitempty"  env:"FLANNELD_ETCD_CERTFILE"`
	EtcdKeyFile   string `yaml:"etcd_keyfile,omitempty"   env:"FLANNELD_ETCD_KEYFILE"`
	EtcdPrefix    string `yaml:"etcd_prefix,omitempty"    env:"FLANNELD_ETCD_PREFIX"`
	IPMasq        string `yaml:"ip_masq,omitempty"        env:"FLANNELD_IP_MASQ"`
	SubnetFile    string `yaml:"subnet_file,omitempty"    env:"FLANNELD_SUBNET_FILE"`
	Iface         string `yaml:"interface,omitempty"      env:"FLANNELD_IFACE"`
	PublicIP      string `yaml:"public_ip,omitempty"      env:"FLANNELD_PUBLIC_IP"`
}

type Fleet struct {
	AgentTTL                string  `yaml:"agent_ttl"                 env:"FLEET_AGENT_TTL"`
	AuthorizedKeysFile      string  `yaml:"authorized_keys_file"      env:"FLEET_AUTHORIZED_KEYS_FILE"`
	DisableEngine           bool    `yaml:"disable_engine"            env:"FLEET_DISABLE_ENGINE"`
	EngineReconcileInterval float64 `yaml:"engine_reconcile_interval" env:"FLEET_ENGINE_RECONCILE_INTERVAL"`
	EtcdCAFile              string  `yaml:"etcd_cafile"               env:"FLEET_ETCD_CAFILE"`
	EtcdCertFile            string  `yaml:"etcd_certfile"             env:"FLEET_ETCD_CERTFILE"`
	EtcdKeyFile             string  `yaml:"etcd_keyfile"              env:"FLEET_ETCD_KEYFILE"`
	EtcdKeyPrefix           string  `yaml:"etcd_key_prefix"           env:"FLEET_ETCD_KEY_PREFIX"`
	EtcdRequestTimeout      float64 `yaml:"etcd_request_timeout"      env:"FLEET_ETCD_REQUEST_TIMEOUT"`
	EtcdServers             string  `yaml:"etcd_servers"              env:"FLEET_ETCD_SERVERS"`
	Metadata                string  `yaml:"metadata"                  env:"FLEET_METADATA"`
	PublicIP                string  `yaml:"public_ip"                 env:"FLEET_PUBLIC_IP"`
	TokenLimit              int     `yaml:"token_limit"               env:"FLEET_TOKEN_LIMIT"`
	Verbosity               int     `yaml:"verbosity"                 env:"FLEET_VERBOSITY"`
	VerifyUnits             bool    `yaml:"verify_units"              env:"FLEET_VERIFY_UNITS"`
}

type Locksmith struct {
	Endpoint           string `yaml:"endpoint"      env:"LOCKSMITHD_ENDPOINT"`
	EtcdCAFile         string `yaml:"etcd_cafile"   env:"LOCKSMITHD_ETCD_CAFILE"`
	EtcdCertFile       string `yaml:"etcd_certfile" env:"LOCKSMITHD_ETCD_CERTFILE"`
	EtcdKeyFile        string `yaml:"etcd_keyfile"  env:"LOCKSMITHD_ETCD_KEYFILE"`
	Group              string `yaml:"group"         env:"LOCKSMITHD_GROUP"`
	RebootWindowStart  string `yaml:"window_start"  env:"REBOOT_WINDOW_START"  valid:"^((?i:sun|mon|tue|wed|thu|fri|sat|sun) )?0*([0-9]|1[0-9]|2[0-3]):0*([0-9]|[1-5][0-9])$"`
	RebootWindowLength string `yaml:"window_length" env:"REBOOT_WINDOW_LENGTH" valid:"^[-+]?([0-9]*(\\.[0-9]*)?[a-z]+)+$"`
}

type OEM struct {
	ID           string `yaml:"id,omitempty"`
	Name         string `yaml:"name,omitempty"`
	VersionID    string `yaml:"version_id,omitempty"`
	HomeURL      string `yaml:"home_url,omitempty"`
	BugReportURL string `yaml:"bug_report_url,omitempty"`
}

type Update struct {
	RebootStrategy string `yaml:"reboot_strategy" env:"REBOOT_STRATEGY" valid:"^(best-effort|etcd-lock|reboot|off)$"`
	Group          string `yaml:"group"           env:"GROUP"`
	Server         string `yaml:"server"          env:"SERVER"`
}

type Unit struct {
	Name    string       `yaml:"name,omitempty"`
	Mask    bool         `yaml:"mask,omitempty"`
	Enable  bool         `yaml:"enable,omitempty"`
	Runtime bool         `yaml:"runtime,omitempty"`
	Content string       `yaml:"content,omitempty"`
	Command string       `yaml:"command,omitempty" valid:"^(start|stop|restart|reload|try-restart|reload-or-restart|reload-or-try-restart)$"`
	DropIns []UnitDropIn `yaml:"drop_ins,omitempty"`
}

type UnitDropIn struct {
	Name    string `yaml:"name,omitempty"`
	Content string `yaml:"content,omitempty"`
}
