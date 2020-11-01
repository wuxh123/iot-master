import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-project-link',
  templateUrl: './project-link.component.html',
  styleUrls: ['./project-link.component.scss']
})
export class ProjectLinkComponent implements OnInit {
  @Input() project: any = {};

  constructor() { }

  ngOnInit(): void {
  }

}
