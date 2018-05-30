// Vis side object
export class GraphVis {
    nodes: VisNode[];
    edges: VisEdge[];
}

export class VisNode {
    id: string;
    label: string;
    group: string;
}

export class VisEdge {
    id: string;
    from: string;
    to: string;
    label: string;
}

// Server side object
export class GraphBean {
    nodes: NodeBean[];
    edges: EdgeBean[];
    options: any;
}

export class NodeBean {
    id: string;
    name: string;
    type: string;
    properties: any;
}

export class EdgeBean {
    id: string;
    from: string;
    to: string;
    type: string;
    properties: any;
}
