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
