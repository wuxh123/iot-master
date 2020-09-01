import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-link-detail',
  templateUrl: './link-detail.component.html',
  styleUrls: ['./link-detail.component.scss']
})
export class LinkDetailComponent implements OnInit {

  @Input() link: any = {};

  constructor() {
  }

  ngOnInit(): void {
  }

}
