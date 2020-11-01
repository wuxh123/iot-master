import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-project-function',
  templateUrl: './project-function.component.html',
  styleUrls: ['./project-function.component.scss']
})
export class ProjectFunctionComponent implements OnInit {
  @Input() project: any = {};

  constructor() { }

  ngOnInit(): void {
  }

}
