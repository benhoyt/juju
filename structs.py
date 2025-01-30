from __future__ import annotations
import dataclasses
from typing import Any


@dataclasses.dataclass
class FormattedBase:
    name: str
    channel: str

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> FormattedBase:
        return cls(
            name=d['name'],
            channel=d['channel'],
        )


@dataclasses.dataclass
class StatusInfoContents:
    current: str | None = None
    message: str | None = None
    reason: str | None = None
    since: str | None = None
    version: str | None = None
    life: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> StatusInfoContents:
        return cls(
            current=d.get('current'),
            message=d.get('message'),
            reason=d.get('reason'),
            since=d.get('since'),
            version=d.get('version'),
            life=d.get('life'),
        )


@dataclasses.dataclass
class AppStatusRelation:
    related_app: str | None = None
    interface: str | None = None
    scope: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> AppStatusRelation:
        return cls(
            related_app=d.get('related-application'),
            interface=d.get('interface'),
            scope=d.get('scope'),
        )


@dataclasses.dataclass
class MeterStatus:
    color: str | None = None
    message: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> MeterStatus:
        return cls(
            color=d.get('color'),
            message=d.get('message'),
        )


@dataclasses.dataclass
class UnitStatus:
    workload_status: StatusInfoContents | None = None
    juju_status: StatusInfoContents | None = None
    meter_status: MeterStatus | None = None
    leader: bool | None = None
    upgrading_from: str | None = None
    machine: str | None = None
    open_ports: list[str] | None = None
    public_address: str | None = None
    address: str | None = None
    provider_id: str | None = None
    subordinates: dict[str, UnitStatus] | None = None
    branch: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> UnitStatus:
        return cls(
            workload_status=StatusInfoContents.from_dict(d['workload-status']) if 'workload-status' in d else None,
            juju_status=StatusInfoContents.from_dict(d['juju-status']) if 'juju-status' in d else None,
            meter_status=MeterStatus.from_dict(d['meter-status']) if 'meter-status' in d else None,
            leader=d.get('leader'),
            upgrading_from=d.get('upgrading-from'),
            machine=d.get('machine'),
            open_ports=d.get('open-ports'),
            public_address=d.get('public-address'),
            address=d.get('address'),
            provider_id=d.get('provider-id'),
            subordinates={k: UnitStatus.from_dict(v) for k, v in d['subordinates'].items()} if 'subordinates' in d else None,
            branch=d.get('branch'),
        )


@dataclasses.dataclass
class AppStatus:
    charm: str
    charm_origin: str
    charm_name: str
    charm_rev: int
    exposed: bool

    base: FormattedBase | None = None
    charm_channel: str | None = None
    charm_version: str | None = None
    charm_profile: str | None = None
    can_upgrade_to: str | None = None
    scale: int | None = None
    provider_id: str | None = None
    address: str | None = None
    life: str | None = None
    app_status: StatusInfoContents | None = None
    relations: dict[str, list[AppStatusRelation]] | None = None
    subordinate_to: list[str] | None = None
    units: dict[str, UnitStatus] | None = None
    version: str | None = None
    endpoint_bindings: dict[str, str] | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> AppStatus:
        return cls(
            charm=d['charm'],
            base=FormattedBase.from_dict(d['base']) if 'base' in d else None,
            charm_origin=d['charm-origin'],
            charm_name=d['charm-name'],
            charm_rev=d['charm-rev'],
            charm_channel=d.get('charm-channel'),
            charm_version=d.get('charm-version'),
            charm_profile=d.get('charm-profile'),
            can_upgrade_to=d.get('can-upgrade-to'),
            scale=d.get('scale'),
            provider_id=d.get('provider-id'),
            address=d.get('address'),
            exposed=d['exposed'],
            life=d.get('life'),
            app_status=StatusInfoContents.from_dict(d['application-status']) if 'application-status' in d else None,
            relations={k: [AppStatusRelation.from_dict(x) for x in v] for k, v in d['relations'].items()} if 'relations' in d else None,
            subordinate_to=d.get('subordinate-to'),
            units={k: UnitStatus.from_dict(v) for k, v in d['units'].items()} if 'units' in d else None,
            version=d.get('version'),
            endpoint_bindings=d.get('endpoint-bindings'),
        )


@dataclasses.dataclass
class BranchStatus:
    ref: str | None = None
    created: str | None = None
    created_by: str | None = None
    active: bool | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> BranchStatus:
        return cls(
            ref=d.get('ref'),
            created=d.get('created'),
            created_by=d.get('created-by'),
            active=d.get('active'),
        )


@dataclasses.dataclass
class EntityStatus:
    current: str | None = None
    message: str | None = None
    since: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> EntityStatus:
        return cls(
            current=d.get('current'),
            message=d.get('message'),
            since=d.get('since'),
        )


@dataclasses.dataclass
class UnitStorageAttachment:
    machine: str | None = None
    location: str | None = None
    life: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> UnitStorageAttachment:
        return cls(
            machine=d.get('machine'),
            location=d.get('location'),
            life=d.get('life'),
        )


@dataclasses.dataclass
class StorageAttachments:
    units: dict[str, UnitStorageAttachment]

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> StorageAttachments:
        return cls(
            units={k: UnitStorageAttachment.from_dict(v) for k, v in d['units'].items()},
        )


@dataclasses.dataclass
class StorageInfo:
    kind: str
    status: EntityStatus
    persistent: bool

    life: str | None = None
    attachments: StorageAttachments | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> StorageInfo:
        return cls(
            kind=d['kind'],
            life=d.get('life'),
            status=EntityStatus.from_dict(d['status']),
            persistent=d['persistent'],
            attachments=StorageAttachments.from_dict(d['attachments']) if 'attachments' in d else None,
        )


@dataclasses.dataclass
class FilesystemAttachment:
    mount_point: str
    read_only: bool

    life: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> FilesystemAttachment:
        return cls(
            mount_point=d['mount-point'],
            read_only=d['read-only'],
            life=d.get('life'),
        )


@dataclasses.dataclass
class FilesystemAttachments:
    machines: dict[str, FilesystemAttachment] | None = None
    containers: dict[str, FilesystemAttachment] | None = None
    units: dict[str, UnitStorageAttachment] | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> FilesystemAttachments:
        return cls(
            machines={k: FilesystemAttachment.from_dict(v) for k, v in d['machines'].items()} if 'machines' in d else None,
            containers={k: FilesystemAttachment.from_dict(v) for k, v in d['containers'].items()} if 'containers' in d else None,
            units={k: UnitStorageAttachment.from_dict(v) for k, v in d['units'].items()} if 'units' in d else None,
        )


@dataclasses.dataclass
class FilesystemInfo:
    Attachments: FilesystemAttachments
    size: int

    provider_id: str | None = None
    volume: str | None = None
    storage: str | None = None
    pool: str | None = None
    life: str | None = None
    status: EntityStatus | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> FilesystemInfo:
        return cls(
            provider_id=d.get('provider-id'),
            volume=d.get('volume'),
            storage=d.get('storage'),
            Attachments=FilesystemAttachments.from_dict(d['Attachments']),
            pool=d.get('pool'),
            size=d['size'],
            life=d.get('life'),
            status=EntityStatus.from_dict(d['status']) if 'status' in d else None,
        )


@dataclasses.dataclass
class VolumeAttachment:
    read_only: bool

    device: str | None = None
    device_link: str | None = None
    bus_address: str | None = None
    life: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> VolumeAttachment:
        return cls(
            device=d.get('device'),
            device_link=d.get('device-link'),
            bus_address=d.get('bus-address'),
            read_only=d['read-only'],
            life=d.get('life'),
        )


@dataclasses.dataclass
class VolumeAttachments:
    machines: dict[str, VolumeAttachment] | None = None
    containers: dict[str, VolumeAttachment] | None = None
    units: dict[str, UnitStorageAttachment] | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> VolumeAttachments:
        return cls(
            machines={k: VolumeAttachment.from_dict(v) for k, v in d['machines'].items()} if 'machines' in d else None,
            containers={k: VolumeAttachment.from_dict(v) for k, v in d['containers'].items()} if 'containers' in d else None,
            units={k: UnitStorageAttachment.from_dict(v) for k, v in d['units'].items()} if 'units' in d else None,
        )


@dataclasses.dataclass
class VolumeInfo:
    size: int
    persistent: bool

    provider_id: str | None = None
    storage: str | None = None
    attachments: VolumeAttachments | None = None
    pool: str | None = None
    hardware_id: str | None = None
    wwn: str | None = None
    life: str | None = None
    status: EntityStatus | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> VolumeInfo:
        return cls(
            provider_id=d.get('provider-id'),
            storage=d.get('storage'),
            attachments=VolumeAttachments.from_dict(d['attachments']) if 'attachments' in d else None,
            pool=d.get('pool'),
            hardware_id=d.get('hardware-id'),
            wwn=d.get('wwn'),
            size=d['size'],
            persistent=d['persistent'],
            life=d.get('life'),
            status=EntityStatus.from_dict(d['status']) if 'status' in d else None,
        )


@dataclasses.dataclass
class CombinedStorage:
    storage: dict[str, StorageInfo] | None = None
    filesystems: dict[str, FilesystemInfo] | None = None
    volumes: dict[str, VolumeInfo] | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> CombinedStorage:
        return cls(
            storage={k: StorageInfo.from_dict(v) for k, v in d['storage'].items()} if 'storage' in d else None,
            filesystems={k: FilesystemInfo.from_dict(v) for k, v in d['filesystems'].items()} if 'filesystems' in d else None,
            volumes={k: VolumeInfo.from_dict(v) for k, v in d['volumes'].items()} if 'volumes' in d else None,
        )


@dataclasses.dataclass
class ControllerStatus:
    timestamp: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> ControllerStatus:
        return cls(
            timestamp=d.get('timestamp'),
        )


@dataclasses.dataclass
class ModelStatus:
    name: str
    type: str
    controller: str
    cloud: str
    version: str

    region: str | None = None
    upgrade_available: str | None = None
    model_status: StatusInfoContents | None = None
    meter_status: MeterStatus | None = None
    sla: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> ModelStatus:
        return cls(
            name=d['name'],
            type=d['type'],
            controller=d['controller'],
            cloud=d['cloud'],
            region=d.get('region'),
            version=d['version'],
            upgrade_available=d.get('upgrade-available'),
            model_status=StatusInfoContents.from_dict(d['model-status']) if 'model-status' in d else None,
            meter_status=MeterStatus.from_dict(d['meter-status']) if 'meter-status' in d else None,
            sla=d.get('sla'),
        )


@dataclasses.dataclass
class NetworkInterface:
    ip_addresses: list[str]
    mac_address: str
    is_up: bool

    gateway: str | None = None
    dns_nameservers: list[str] | None = None
    space: str | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> NetworkInterface:
        return cls(
            ip_addresses=d['ip-addresses'],
            mac_address=d['mac-address'],
            gateway=d.get('gateway'),
            dns_nameservers=d.get('dns-nameservers'),
            space=d.get('space'),
            is_up=d['is-up'],
        )


@dataclasses.dataclass
class LxdProfileContents:
    config: dict[str, str]
    description: str
    devices: dict[str, dict[str, str]]

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> LxdProfileContents:
        return cls(
            config=d['config'],
            description=d['description'],
            devices=d['devices'],
        )


@dataclasses.dataclass
class MachineStatus:
    juju_status: StatusInfoContents | None = None
    hostname: str | None = None
    dns_name: str | None = None
    ip_addresses: list[str] | None = None
    instance_id: str | None = None
    display_name: str | None = None
    machine_status: StatusInfoContents | None = None
    modification_status: StatusInfoContents | None = None
    base: FormattedBase | None = None
    network_interfaces: dict[str, NetworkInterface] | None = None
    containers: dict[str, MachineStatus] | None = None
    constraints: str | None = None
    hardware: str | None = None
    controller_member_status: str | None = None
    ha_primary: bool | None = None
    lxd_profiles: dict[str, LxdProfileContents] | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> MachineStatus:
        return cls(
            juju_status=StatusInfoContents.from_dict(d['juju-status']) if 'juju-status' in d else None,
            hostname=d.get('hostname'),
            dns_name=d.get('dns-name'),
            ip_addresses=d.get('ip-addresses'),
            instance_id=d.get('instance-id'),
            display_name=d.get('display-name'),
            machine_status=StatusInfoContents.from_dict(d['machine-status']) if 'machine-status' in d else None,
            modification_status=StatusInfoContents.from_dict(d['modification-status']) if 'modification-status' in d else None,
            base=FormattedBase.from_dict(d['base']) if 'base' in d else None,
            network_interfaces={k: NetworkInterface.from_dict(v) for k, v in d['network-interfaces'].items()} if 'network-interfaces' in d else None,
            containers={k: MachineStatus.from_dict(v) for k, v in d['containers'].items()} if 'containers' in d else None,
            constraints=d.get('constraints'),
            hardware=d.get('hardware'),
            controller_member_status=d.get('controller-member-status'),
            ha_primary=d.get('ha-primary'),
            lxd_profiles={k: LxdProfileContents.from_dict(v) for k, v in d['lxd-profiles'].items()} if 'lxd-profiles' in d else None,
        )


@dataclasses.dataclass
class RemoteEndpoint:
    interface: str
    role: str

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> RemoteEndpoint:
        return cls(
            interface=d['interface'],
            role=d['role'],
        )


@dataclasses.dataclass
class RemoteAppStatus:
    url: str

    endpoints: dict[str, RemoteEndpoint] | None = None
    life: str | None = None
    app_status: StatusInfoContents | None = None
    relations: dict[str, list[str]] | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> RemoteAppStatus:
        return cls(
            url=d['url'],
            endpoints={k: RemoteEndpoint.from_dict(v) for k, v in d['endpoints'].items()} if 'endpoints' in d else None,
            life=d.get('life'),
            app_status=StatusInfoContents.from_dict(d['application-status']) if 'application-status' in d else None,
            relations=d.get('relations'),
        )


@dataclasses.dataclass
class OfferStatus:
    app: str
    endpoints: dict[str, RemoteEndpoint]

    charm: str | None = None
    total_connected_count: int | None = None
    active_connected_count: int | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> OfferStatus:
        return cls(
            app=d['application'],
            charm=d.get('charm'),
            total_connected_count=d.get('total-connected-count'),
            active_connected_count=d.get('active-connected-count'),
            endpoints={k: RemoteEndpoint.from_dict(v) for k, v in d['endpoints'].items()},
        )


@dataclasses.dataclass
class FormattedStatus:
    model: ModelStatus
    machines: dict[str, MachineStatus]
    apps: dict[str, AppStatus]

    app_endpoints: dict[str, RemoteAppStatus] | None = None
    offers: dict[str, OfferStatus] | None = None
    storage: CombinedStorage | None = None
    controller: ControllerStatus | None = None
    branches: dict[str, BranchStatus] | None = None

    @classmethod
    def from_dict(cls, d: dict[str, Any]) -> FormattedStatus:
        return cls(
            model=ModelStatus.from_dict(d['model']),
            machines={k: MachineStatus.from_dict(v) for k, v in d['machines'].items()},
            apps={k: AppStatus.from_dict(v) for k, v in d['applications'].items()},
            app_endpoints={k: RemoteAppStatus.from_dict(v) for k, v in d['application-endpoints'].items()} if 'application-endpoints' in d else None,
            offers={k: OfferStatus.from_dict(v) for k, v in d['offers'].items()} if 'offers' in d else None,
            storage=CombinedStorage.from_dict(d['storage']) if 'storage' in d else None,
            controller=ControllerStatus.from_dict(d['controller']) if 'controller' in d else None,
            branches={k: BranchStatus.from_dict(v) for k, v in d['branches'].items()} if 'branches' in d else None,
        )
