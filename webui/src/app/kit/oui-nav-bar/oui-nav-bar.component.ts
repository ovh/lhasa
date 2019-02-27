import { Component, OnInit, Output, EventEmitter, Input, ElementRef, OnDestroy, Renderer2 } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { ApplicationBean } from '../../models/commons/applications-bean';
import { Store } from '@ngrx/store';
import { ApplicationsStoreService } from '../../stores/applications-store.service';
import { UiKitMenuItem } from '../../models/kit/navbar';
import { DomHandler } from 'primeng/primeng';
import { each } from 'lodash';
import { HelpsStoreService } from '../../stores/help-store.service';

@Component({
  selector: 'app-oui-nav-bar',
  templateUrl: './oui-nav-bar.component.html',
  styleUrls: ['./oui-nav-bar.component.css'],
  providers: [DomHandler]
})
export class OuiNavBarComponent implements OnInit {

  visible = true;
  @Input() home = '/';
  @Input() items: UiKitMenuItem;

  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
    private el: ElementRef,
    private domHandler: DomHandler,
    private renderer: Renderer2,
    public help: HelpsStoreService
  ) {
  }

  ngOnInit() {
  }

  onSelect(event: any, tabs: string) {
    event.data = tabs;
    this.select.emit(event);
  }

  release() {
    each(this.items, (item) => {
      each(item.items, (subitem) => {
        each(DomHandler.find(this.el.nativeElement, '#button-' + subitem.id), (elem) => {
          elem.expanded = false;
          elem.setAttribute('aria-expanded', 'false');
        });
      });
      each(DomHandler.find(this.el.nativeElement, '#button-' + item.id), (elem) => {
        elem.expanded = false;
        elem.setAttribute('aria-expanded', 'false');
      });
    });
  }

  expand(item: UiKitMenuItem) {
    item.expanded = !item.expanded;
    const elem = DomHandler.findSingle(this.el.nativeElement, '#button-' + item.id);
    elem.setAttribute('aria-expanded', String(item.expanded));
  }
}
