export class EntityBean {
    id: string
    timestamp: Date
}

export class ContentListResponse<T extends EntityBean> {
    _links: {}
    start: number;
    size: number;
    limit: number;
    content: T[];
}
