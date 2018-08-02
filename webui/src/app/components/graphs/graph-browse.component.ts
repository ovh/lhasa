import { Component, OnInit, AfterViewInit, ViewChild } from '@angular/core';
import { GraphBean, NodeBean, GraphVis } from '../../models/graph/graph-bean';
import { Observable } from 'rxjs';
import { AutoUnsubscribe } from '../../shared/decorator/autoUnsubscribe';
import { GraphsStoreService } from './../../stores/graphs-store.service';
import { CytoscapeModule, CytoscapeComponent } from 'ngx-cytoscape';
import { environment } from '../../../environments/environment';
import { Router } from '@angular/router';
const cytoscape = require('cytoscape');
const cytoscapeExpandCollapse = require('cytoscape-expand-collapse');
const coseBilkent = require('cytoscape-cose-bilkent');
const jquery = require('jquery');

@Component({
  selector: 'app-graph-browse',
  templateUrl: './graph-browse.component.html',
  styleUrls: ['./graph-browse.component.css'],
})

@AutoUnsubscribe()
export class GraphBrowseComponent implements OnInit, AfterViewInit {

  @ViewChild('cytograph') cytograph: CytoscapeComponent;

  public layout: any = {};


  public style: any = [
    {
      selector: 'node',
      style: {
        'background-color': 'gray',
        'label': 'data(name)',
        'font-size': '30',
      }
    },
    {
      selector: 'node.environment, node.domain',
      style: {
        'shape': 'square',
        'background-color': 'data(color)',
        'color': 'data(color)',
        'font-size': '30',
        // 'size': '300',
      }
    },
    {
      selector: '$node > node',
      css: {
        'padding-top': '10px',
        'padding-left': '10px',
        'padding-bottom': '10px',
        'padding-right': '10px',
        'text-valign': 'top',
        'text-halign': 'center',
        'background-color': 'white',
        'border-width': '5px',
        'border-color': 'data(color)',
        'font-size': '100',
        'text-outline-width': 10,
        'text-outline-color': 'data(color)',
        'color': '#fff'
      }
    },
    {
      selector: 'edge',
      style: {
        'width': 0.5,
        'line-color': '#bbb',
        'curve-style': 'unbundled-bezier',
        'target-arrow-color': '#bbb',
        'target-arrow-shape': 'triangle',
      }
    },
    {
      selector: 'node.semitransp',
      style: {
        'opacity': '0.1',
      }
    },
    {
      selector: 'node.incomer',
      style: {
        'border-color': 'orange',
        'border-width': '10px',
        'opacity': '1',
      }
    },
    {
      selector: 'node.outgoer',
      style: {
        'border-color': 'purple',
        'border-width': '10px',
        'opacity': '1',
      }
    },
    {
      selector: 'edge.outgoer',
      style: {
        'line-color': 'purple',
        'opacity': '0.3',
        'target-arrow-color': 'purple',
        'width': 10,
      }
    },
    {
      selector: 'edge.incomer',
      style: {
        'line-color': 'orange',
        'opacity': '0.3',
        'target-arrow-color': 'orange',
        'width': 10,
      }
    },
  ];

  private _graphData: any = {
    nodes: [],
    edges: []
  };
  protected deploymentStream: Observable<GraphBean>;
  public deployments: GraphBean = {
    nodes: [],
    edges: [],
    options: {}
  };

  public deploymentsVis: GraphVis = {
    nodes: [],
    edges: [],
  };

  private api: any;

  constructor(
    private router: Router,
    private graphsStoreService: GraphsStoreService,
  ) {
    this.deploymentStream = this.graphsStoreService.deployments();
  }



  ngOnInit() {
    this.deploymentStream.subscribe(
      (graph: GraphBean) => {
        this.deployments = graph;
        const dejaVuDomains = {};
        const dejaVuEnvironments = {};
        const dejaVuNodes = {};
        if (graph.nodes === undefined || graph.edges === undefined) {
          return;
        }
        graph.nodes.forEach((node, index) => {
          if (index > 10000) {
            return;
          }
          dejaVuNodes[node.id] = true;
          const env = node.properties.environment.slug;
          const domain = env + '/' + node.properties.application.domain;
          // const domain = node.properties.application.domain;
          if (dejaVuEnvironments[env] === undefined) {
            this._graphData.nodes.push({
              classes: 'environment',
              data: {
                id: env,
                type: 'environment',
                color: node.properties.environment.properties.color  || 'gray',
                name: env,
                // faveShape: 'triangle'
              }
            });
            dejaVuEnvironments[env] = true;
          }
          if (dejaVuDomains[domain] === undefined) {
            this._graphData.nodes.push({
              classes: 'domain',
              data: {
                id: domain,
                name: node.properties.application.domain,
                color: node.properties.environment.properties.color || 'gray',
                type: 'domain',
                parent: env,
              }
            });
            dejaVuDomains[domain] = true;
          }
          this._graphData.nodes.push({
            classes: 'application',
            data: {
              id: node.id,
              name: node.name,
              type: 'application',
              domain: node.properties.application.domain,
              parent: domain,
            }
          });
        });
        graph.edges.forEach(edge => {
          if (dejaVuNodes[edge.from] === undefined || dejaVuNodes[edge.to] === undefined) {
            return;
          }
          this._graphData.edges.push({
            data: {
              source: edge.from,
              target: edge.to,
              // faveColor: 'lightgray'
            }
          });
        });
      },
      error => {
        console.error(error);
      },
      () => {
      }
    );
  }

  ngAfterViewInit() {
    // cytoscapeExpandCollapse(cytoscape, jquery);
    coseBilkent(window['cytoscape']);
    window['cytoscapeExpandCollapse'](window['cytoscape'], jquery);
    this.api = this.cytograph.cy.expandCollapse({
      layoutBy: {
        name: 'cose-bilkent',
        nodeDimensionsIncludeLabels: true,
        refresh: 250,
        animate: false,
        randomize: false,
        fit: true,
        nodeRepulsion: 10e6,
        idealEdgeLength: 100,
        edgeElasticity: 1,
        nestingFactor: 0.001,
      },
      undoable: false,
      fisheye: true,
    });
    window['api'] = this.api;

    const cy = this.cytograph.cy;
    const component = this;
    cy.on('mouseover', 'node', (e) => {
      this.api.setOption('layoutBy', {name: 'preset'});
      const sel = e.target;
      if (sel.isParent() && ! sel.hasClass('cy-expand-collapse-collapsed-node')) { return; }
      const related = sel.outgoers().addClass('outgoer').union(sel.incomers().addClass('incomer'));
      related.merge(related.ancestors()).merge(sel);
      e.cy.elements().difference(related).difference(sel).addClass('semitransp');
    });
    cy.on('mouseout', 'node', (e) => {
      const sel = e.target;
      e.cy.elements().removeClass('outgoer').removeClass('incomer').removeClass('semitransp');
    });
    cy.on('tap', 'node.application', function() {
      const app = this.data('name');
      const domain = this.data('domain');
      const url = `/applications/${domain}/${app}`;
      if (confirm('Are you sure you want to leave the map ? (all changes will be lost)')) {
        component.router.navigateByUrl(url);
      }
    });
    this.api.collapse(cy.elements('.domain'));
  }
  get graphData(): any {
    return this._graphData;
  }

  set graphData(value: any) {
    this._graphData = value;
  }
}

