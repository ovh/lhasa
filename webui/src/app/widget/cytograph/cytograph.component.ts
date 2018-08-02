import { Component, OnInit, ViewChild } from '@angular/core';
import { AfterViewInit } from '@angular/core/src/metadata/lifecycle_hooks';
import { CytoscapeComponent } from 'ngx-cytoscape';
import { Router } from '@angular/router';

const cytoscape = require('cytoscape');
const cytoscapeExpandCollapse = require('cytoscape-expand-collapse');
const coseBilkent = require('cytoscape-cose-bilkent');
const jquery = require('jquery');

@Component({
  selector: 'app-cytograph',
  templateUrl: './cytograph.component.html',
  styleUrls: ['./cytograph.component.css'],
})

export class CytoGraphComponent implements OnInit, AfterViewInit {
  @ViewChild('cytograph') cytograph: CytoscapeComponent;
  public graphData: any;
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

  private cytoscapeExpandCollapse: any;

  constructor(private router: Router) {
  }

  ngOnInit() {
  }

  load(graphData: any) {
    this.graphData = graphData;
  }

  ngAfterViewInit() {
    coseBilkent(cytoscape);
    cytoscapeExpandCollapse(window['cytoscape'], jquery);
    this.cytoscapeExpandCollapse = this.cytograph.cy.expandCollapse({
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

    const cy = this.cytograph.cy;
    const component = this;

    cy.on('mouseover', 'node', (e) => {
      this.cytoscapeExpandCollapse.setOption('layoutBy', {name: 'preset'});
      const sel = e.target;
      if (sel.isParent() && ! sel.hasClass('cy-expand-collapse-collapsed-node')) {
        return;
      }
      const related = sel.outgoers().addClass('outgoer').union(sel.incomers().addClass('incomer'));
      related.merge(sel);
      related.merge(related.ancestors());
      e.cy.elements().difference(related).addClass('semitransp');
    });

    cy.on('mouseout', 'node', (e) => {
      const sel = e.target;
      e.cy.elements().removeClass('outgoer').removeClass('incomer').removeClass('semitransp');
    });

    this.cytoscapeExpandCollapse.collapse(cy.elements('.domain'));
  }
}
