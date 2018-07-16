import { EntityBean, PageMetaData, HrefLinks, AbstractPaginatedResource } from './entity-bean';
import { Timestamp } from 'rxjs/operators/timestamp';

export class BadgeLevelBean {
  id: string;
  description: string;
  label: string;
  color: string;
}

// Badge
export class BadgeBean extends EntityBean {
  slug: string;
  title: string;
  type: string;
  levels: BadgeLevelBean[];
  _links?: HrefLinks[];
  _stats: Map<string, number>;
}

// Badge for page browse
export class BadgePagesBean  extends AbstractPaginatedResource {
  badges: BadgeBean[] = [];
  metadata: PageMetaData = {
    totalElements: 0,
    totalPages: 0,
    size: 0,
    number: 0
  };
}
