export class EntityBean {
    id?: string;
    timestamp?: Date;
}

export class HrefLinks {
    rel: string;
    href: string;
}

// Page meta data
export class PageMetaData {
    totalElements?: number;
    totalPages?: number;
    size: number;
    number: number;
}

export class ContentListResponse<T extends EntityBean> {
    content: T[];
    pageMetadata: PageMetaData;
    _links: {};
}


export abstract class AbstractPaginatedResource {
    metadata: PageMetaData = {
      totalElements: 0,
      totalPages: 0,
      size: 0,
      number: 0
    };
  }
