import { Component, OnInit, Output, EventEmitter, Input, ElementRef, OnDestroy, Renderer2 } from '@angular/core';
import { DomHandler } from 'primeng/primeng';

// BadgeShieldsIOBean for a badge SVG representation
export class BadgeShieldsIOBean {
  title: string;
  value: string;
  comment: string;
  label: string;
  color: string;
  description: string;
}

@Component({
  selector: 'app-badge-shieldsio',
  templateUrl: './badgewidget.component.html',
  styleUrls: ['./badgewidget.component.css'],
  providers: [DomHandler]
})
export class BadgeWidgetComponent implements OnInit {

  _badge: BadgeShieldsIOBean;

  @Output() select: EventEmitter<any> = new EventEmitter();

  constructor(
  ) {
  }

  ngOnInit() {
  }

  @Input() set badge(val: BadgeShieldsIOBean) {
    this._badge = val;
  }
}
