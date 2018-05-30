import { DomHandler } from 'primeng/primeng';
import { Component, Input, OnInit, ElementRef, ViewChild, Output, EventEmitter } from '@angular/core';
import { DataSet, IdType, Network } from 'vis';
import { AfterViewInit } from '@angular/core/src/metadata/lifecycle_hooks';
import { GraphVis } from '../../models/graph/graph-bean';
import { DataContentService } from '../../services/data-content.service';
import { ContentBean } from '../../models/commons/content-bean';

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
        // Not initialized
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
            edges: this._graph.edges
        };

        this.contentService.GetSingle(this._options).subscribe(
            (content: ContentBean) => {
                this.netOptions = JSON.parse(<any>content);

                if (this.network) {
                    // Cf. http://visjs.org/docs/network for documentation
                    this.network.setData(data.nodes, data.edges);
                } else {
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
                this.network.redraw();
                this.network.setSize(this.container.nativeElement.width, 800);
            }
        )
    }
}
