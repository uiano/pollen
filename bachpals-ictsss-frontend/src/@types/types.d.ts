export interface SERVER_INFO {
  ServerIp: string;
  ServerImage: string;
  ServerName: string;
  UserId: string;
  ServerId: string;
  Created: string;
  ServerStatus: OS_STATUS;
  GroupMembers: Array<string>;
  ImageReadRootPassword: string;
  ImageDisplayName: string;
}

export type OS_STATUS = "ACTIVE" | "SHUTOFF";

export type VMS_ARRAY = Array<SERVER_INFO>;

export type VMS_RESPONSE = { data: VMS_ARRAY };

export type ITabContext = {
  selectedTab: number;
  setSelectedTab: React.Dispatch<React.SetStateAction<number>> | null;
};

export type Administrators = {
  Name: string;
  UserId: string;
};

export type Image = {
  Id: string;
  Published: string;
  ImageId: string;
  ImageName: string;
  ImageDescription: string;
  ImageDisplayName: string;
  ImageConfig: string;
  ImageReadRootPassword: boolean;
};

export type ImageConfig = string;

export type ServerImage = {
  ImageId: string;
  Name: string;
};

export type SelectServerImage = {
  ImageId: string;
  ImageDisplayName: string;
};

export type AuthContextProps = {
  user: User | null;
  isLoading: boolean;
};

export type User = {
  "connect-userid_sec"?: Array<string>;
  "dataporten-userid_sec"?: Array<string>;
  email?: string;
  email_verified?: boolean;
  name?: string;
  sub?: string;
  token?: string;
};
