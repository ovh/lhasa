import { Component, OnInit, Output, EventEmitter, Input, OnDestroy } from '@angular/core';

import { UiKitMenuItem } from '../../models/kit/navbar';

@Component({
  selector: 'app-oui-action-menu',
  templateUrl: './oui-action-menu.component.html',
  styleUrls: ['./oui-action-menu.component.css'],
  providers: []
})
export class OuiActionMenuComponent implements OnInit {

  visible = true;
  @Input() items: UiKitMenuItem;

  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
  ) {
  }

  ngOnInit() {
  }

  onSelect(event: any, tabs: string) {
    event.data = tabs;
    this.select.emit(event);
  }
}
