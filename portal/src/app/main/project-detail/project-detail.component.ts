import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {TabRef} from '../tabs/tabs.component';
import {ApiService} from '../../api.service';

@Component({
  selector: 'app-project-detail',
  templateUrl: './project-detail.component.html',
  styleUrls: ['./project-detail.component.scss']
})
export class ProjectDetailComponent implements OnInit {

  project: any = {};
  id = 0;

  constructor(private as: ApiService, private routeInfo: ActivatedRoute, private tab: TabRef) {
    this.id = routeInfo.snapshot.params.id;
    tab.name = '项目详情';
  }

  ngOnInit(): void {
    this.as.get('project/' + this.id).subscribe(res => {
      if (res.ok) {
        this.project = res.data;
        this.tab.name = '项目【' + this.project.name + '】';
      }
      // console.log(res);
    });
  }

}
