import { DomHandler } from 'primeng/primeng';
import { Component, Input, OnInit, ElementRef, ViewChild, Output, EventEmitter } from '@angular/core';
import { DataSet, IdType, Network } from 'vis';
import { AfterViewInit } from '@angular/core/src/metadata/lifecycle_hooks';
import { GraphVis } from '../../models/graph/graph-bean';
import { DataContentService } from '../../services/data-content.service';
import { ContentBean } from '../../models/commons/content-bean';
import { environment } from '../../../environments/environment.prod';
import { MessagesModule } from 'primeng/messages';
import { Message } from 'primeng/api';

@Component({
    selector: 'app-graph',
    templateUrl: './graph.component.html',
    styleUrls: ['./graph.component.css']
})

export class GraphComponent implements OnInit, AfterViewInit {

    /**
     * internal members
     */
    public display = false;
    protected _graph: GraphVis;
    protected _options: any;
    protected netOptions: any;
    protected network: any;
    public msgs: Message[] = [];


    @ViewChild('graphvis') container: any;
    @ViewChild('graphvisconfig') config: any;
    @Output() handler: EventEmitter<any> = new EventEmitter();

    /**
     * internal
     */

    constructor(
        private elementRef: ElementRef,
        private contentService: DataContentService,
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
        this.update();
    }

    @Input() get graph(): GraphVis {
        return this._graph;
    }

    set graph(val: GraphVis) {
        this._graph = val;
    }

    @Input() get options(): any {
        return this._options;
    }

    set options(val: any) {
        this._options = val;
    }

    /**
     * configure network
     */
    public configure() {
        if (this.netOptions.configure && this.netOptions.configure.container) {
            this.display = true;
            return;
        }
        this.netOptions.configure = {};
        this.netOptions.configure.enabled = true;
        this.netOptions.configure.showButton = true;
        this.netOptions.configure.container = this.config.nativeElement;
        this.network.setOptions(this.netOptions);
        this.display = true;
    }
    /**
     * update edge
     */
    public update() {
        const data = {
            nodes: this._graph.nodes,
            edges: this._graph.edges,
        };

        // HACK: This removes the edges that start and end at the same node (loop)
        // Vis.js does not  like clustering + edge loops :(
        data.edges = data.edges.filter(edge => edge.from !== edge.to);

        // HACK: This limits the edge count on Fireox and displays a warning message
        var edgeCountLimit = 100;
        if (data.edges.length > edgeCountLimit) {
            var msg = `This graph may take up-to 2 minutes to display due to its size and to some performance issues of the vis.js plugin. We're working to fix that.`
            if (navigator.userAgent.toLowerCase().indexOf('firefox') > -1) {
                var totalLen = data.edges.length;
                msg = `This graph is incomplete. It will only display the first ${edgeCountLimit} edges out of ${totalLen} due to a performance limitation of the vis.js plugin with Firefox. Please use Chromium for now if you want to display the whole graph. We're working to fix that.`
                data.edges = data.edges.slice(0, edgeCountLimit)
            }
            this.msgs.push({
                severity: 'warn',
                summary: 'Experimental feature',
                detail: msg,
            });
        }

        this.contentService.GetSingle(this._options).subscribe(
            async (content: ContentBean) => {
                // Hack: This delay allows to let the time to display the warning message
                await delay(500);
                this.netOptions = {}
                if (content) {
                    this.netOptions = JSON.parse(<any>content);
                }
                else {

                    this.netOptions = {
                        autoResize: true,
                        edges: {
                            arrows: {
                                to: {
                                    enabled: true,
                                }
                            }
                        },
                        nodes: {
                            shadow: {
                                enabled: true,
                            },
                            shape: 'dot',
                            color: {
                                border: '#2B7CE9',
                                background: '#122844',
                            },
                        },
                        physics: {
                            maxVelocity: 3,
                        }
                    }
                }
                if (this.network) {
                    // Cf. http://visjs.org/docs/network for documentation
                    this.network.setData(data.nodes, data.edges);
                } else {
                    var domains = new Set();
                    data.nodes.forEach(node => {
                        domains.add(node['domain'])
                    })
                    // create a network
                    this.network = new Network(this.container.nativeElement, data, this.netOptions);
                    [
                        'click',
                        'doubleClick',
                        'oncontext',
                        'dragStart',
                        'dragging',
                        'dragEnd',
                        'zoom',
                        'showPopup',
                        'hidePopup',
                        'showPopup',
                        'select',
                        'selectNode',
                        'selectEdge',
                        'deselectNode',
                        'deselectEdge',
                        'hoverNode',
                        'hoverEdge',
                        'blurNode',
                        'blurEdge'
                    ].forEach(name => {
                        this.network.on(name, (params) => {
                            params.type = name;
                            this.handler.emit(params);
                        });
                    });
                }
                var i = 1;
                domains.forEach(domain => {
                    var clusterOptionsByData = {
                        joinCondition: function (childOptions) {
                            return childOptions.domain == domain;
                        },
                        clusterNodeProperties: {
                            id: 'domain:' + domain,
                            borderWidth: 3,
                            shape: 'dot',
                            label: domain,
                            allowSingleNodeCluster: false
                        }
                    };
                    this.network.cluster(clusterOptionsByData);
                    console.log(`clustering domain ${domain} ${i}/${domains.size}`);
                    i++;
                })
                const net = this.network;
                this.network.on("selectNode", function (params) {
                    if (params.nodes.length == 1) {
                        if (net.isCluster(params.nodes[0]) == true) {
                            net.openCluster(params.nodes[0]);
                        }
                    }
                });
                this.network.setSize(this.container.nativeElement.width, "1000");
                this.network.redraw();
            }
        )
    }
}

function delay(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}