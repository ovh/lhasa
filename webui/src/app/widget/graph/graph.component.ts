import {Component, Input, OnInit} from '@angular/core';
import {DataSet, IdType, Network} from 'vis';
import {AfterViewInit} from '@angular/core/src/metadata/lifecycle_hooks';
import {Edge, Node} from '../../models/graph/graph-bean';

@Component({
    selector: 'app-graph',
    templateUrl: './graph.component.html',
    styleUrls: ['./graph.component.css']
})
export class GraphComponent implements OnInit, AfterViewInit {

    protected _nodes: Node[];
    protected _edges: Edge[];

    /**
     * internal
     */

    constructor(
    ) {
    }

    /**
     * init component
     */
    ngOnInit() {
    }

    /**
     * init component
     */
    ngAfterViewInit() {
        setTimeout(() => {
            this.update();
        }, 1000);
    }

    @Input() get nodes(): Node[] {
        return this._nodes;
    }

    set nodes(val: Node[]) {
        this._nodes = val;
    }

    @Input() get edges(): Edge[] {
        return this._edges;
    }

    set edges(val: Edge[]) {
        this._edges = val;
    }

    /**
     * update edge
     */
    public update() {
        // create a network
        const container = document.getElementById('mynetwork');
        const data = {
            nodes: this._nodes,
            edges: this._edges
        };
        const options = {
            groups: {
                failure: {
                    color: {
                        background: 'red'
                    }
                },
                state: {
                    color: {
                        background: 'lime'
                    }
                },
                startstate: {
                    font: {
                        size: 33,
                        color: 'white'
                    },
                    color: {
                        background: 'blueviolet'
                    }
                },
                finalstate: {
                    font: {
                        size: 33,
                        color: 'white'
                    },
                    color: {
                        background: 'blue'
                    }
                }
            },
            edges: {
                arrows: {
                    to: {
                        enabled: true
                    }
                },
                smooth: false
            },
            physics: {
                adaptiveTimestep: true,
                barnesHut: {
                    gravitationalConstant: -8000,
                    springConstant: 0.04,
                    springLength: 95
                },
                stabilization: false
            },
            layout: {
                randomSeed: 191006,
                improvedLayout: false
            }
        };
        var network = new Network(container, data, options);
    }

}
