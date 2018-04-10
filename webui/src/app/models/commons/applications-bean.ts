import { EntityBean } from "./entity-bean";

// Application
export class ApplicationBean extends EntityBean {
  domain: string;
  name: string;
  type: string;
  language: string;
  repositoryurl: string;
  project: string;
  repo: string;
  avatarurl: string;
  description: string;
  manifest: ManifestBean
}

// Manifest
export class ManifestBean {
  name: string
  profile: string
  description: string
  repository: string
  authors: PersonBean[]
  support: TeamBean
}

// Person
export class PersonBean {
  name: string;
  email: string;
  cisco: string;
}

// Team
export class TeamBean {
  name: string;
  email: string;
  cisco: string;
}

// Domain
export class DomainBean {
  name: string;
  applications: ApplicationBean[];
}

// Bitbucket
export class BitbucketBean extends EntityBean {
}
