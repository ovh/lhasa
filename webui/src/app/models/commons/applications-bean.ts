import { EntityBean } from "./entity-bean";

export class ApplicationBean extends EntityBean {
    domain: string;
    name: string;
    type: string;
    language: string;
    repositoryurl: string;
    avatarurl: string;
    description: string;
    manifest: ManifestBean
  }

  export class ManifestBean {
    description: string;
  }
